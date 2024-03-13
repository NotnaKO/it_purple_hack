package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

var metrics *Metrics

func NewMetrics() *Metrics {
	return &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "price_manager_requests_total",
				Help: "Total number of requests processed by the Price Manager service.",
			},
			[]string{"method", "path"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "price_manager_request_duration_seconds",
				Help: "Duration of requests processed by the Price Manager service.",
			},
			[]string{"method", "path"},
		),
	}
}

func (m *Metrics) Register() {
	prometheus.MustRegister(m.RequestsTotal)
	prometheus.MustRegister(m.RequestDuration)
}

func InitializeMetrics() {
	metrics = NewMetrics()
	metrics.Register()
	http.Handle("/metrics", promhttp.Handler())
}

type JSONCategory struct {
	Name string `json:"name"`
	Id   uint64 `json:"id"`
}

type Handler struct {
	logger       *logrus.Logger
	priceManager *PriceManager
}

func NewHandler(priceManager *PriceManager, logger *logrus.Logger) *Handler {
	return &Handler{
		logger:       logger,
		priceManager: priceManager,
	}
}

func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	metrics.RequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()

	startTime := time.Now()

	getRequest, err := NewGetRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, err := h.priceManager.GetPrice(&getRequest)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Price uint64 `json:"price"`
	}{
		Price: price,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logRequest(r, startTime)
}

func (h *Handler) SetPrice(w http.ResponseWriter, r *http.Request) {
	metrics.RequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()

	startTime := time.Now()

	setRequest, err := NewSetRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.priceManager.SetPrice(&setRequest)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Log request details
	h.logRequest(r, startTime)
}

func (h *Handler) logRequest(r *http.Request, startTime time.Time) {
	duration := time.Since(startTime)
	metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
	h.logger.WithFields(logrus.Fields{
		"method":   r.Method,
		"path":     r.URL.Path,
		"duration": duration.Seconds(),
	}).Info("Request processed")
}

func (h *Handler) logRequestError(r *http.Request, err error) {
	h.logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"error":  err.Error(),
	}).Error("Request error")
}

func (h *Handler) logServerError(r *http.Request, err error) {
	h.logger.WithFields(logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"error":  err.Error(),
	}).Error("Internal Server Error")
}

func ConnectToDatabase() (*sql.DB, error) {
	logrus.Info("Connect to database:", fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable", config.PostgresqlUser,
		config.Password, config.PostgresqlHost, config.Dbname))
	return sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.PostgresqlUser, config.Password, config.PostgresqlHost, config.Dbname))
}

var configPath = flag.String("config_path", "",
	"Path to the retrieval file .yaml file which contains all field of config. db_path_name should be json file")
var config Config
var NoConfig = errors.New("you should set config file. Use --help to information")

func ParseTableIdJson(filename string) (map[uint64]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(file)

	var tables []JSONCategory
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tables); err != nil {
		return nil, err
	}
	DataBaseById := make(map[uint64]string)
	for _, table := range tables {
		DataBaseById[table.Id] = table.Name
	}
	return DataBaseById, nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	flag.Parse()
	if *configPath == "" {
		logrus.Fatal(NoConfig)
	}
	err := error(nil)
	config, err = loadConfig(*configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	db, err := ConnectToDatabase()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logrus.Fatal(err)
		}
	}(db)

	priceManager, err := NewPriceManagementService(db, config.DbPathName)
	if err != nil {
		logrus.Fatal(err)
	}

	err = priceManager.loadDB()
	if err != nil {
		logrus.Fatal(err)
	}

	// Now we are ready to start
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel) // Set log level to debug

	InitializeMetrics()
	logrus.Info("Prometheus initialized successfully")

	handler := NewHandler(priceManager, logger)
	http.HandleFunc("/get_price", handler.GetPrice)
	http.HandleFunc("/set_price", handler.SetPrice)
	// TODO kubernetes. right now leave only one port
	go func() {
		port := strconv.Itoa(int(config.ServerPort))
		fmt.Printf("Price Management Service is listening on port %s...\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	select {}
}
