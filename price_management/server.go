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

func (h *Handler) GetMatrixById(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	getRequest, err := NewGetMatrixByIdRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matrix, err := h.priceManager.GetMatrixById(&getRequest)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Matrix string `json:"matrix_name"`
	}{
		Matrix: matrix,
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

func (h *Handler) GetIdByMatrix(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	getRequest, err := NewGetIdByMatrixRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.priceManager.GetIdByMatrix(&getRequest)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Matrix uint64 `json:"id_matrix"`
	}{
		Matrix: id,
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

func (h *Handler) ChangeStorage(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	getRequest, err := NewChangeStorage(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.priceManager.ChangeStorage(&getRequest)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%b", id)

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

func (p *PriceManager) ParseTableIdJson(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
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
		return err
	}
	p.DataBaseById = make(map[uint64]string)
	for _, table := range tables {
		p.DataBaseById[table.Id] = table.Name
	}
	return nil
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
	logrus.Debug("Connect successful")

	priceManager, err := NewPriceManagementService(db, config.DbPathName)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debug("New price manager created")

	err = priceManager.loadDB()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Debug("Db load successfully")

	defer func(priceManager *PriceManager) {
		err := priceManager.dumpTables()
		if err != nil {
			logrus.Error("Couldn't dump tables in defer, because get error: ", err)
		}
	}(priceManager)

	tickChannel := time.Tick(time.Minute)
	go priceManager.waitTimer(tickChannel)
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
	http.HandleFunc("/get_matrix", handler.GetMatrixById)
	http.HandleFunc("/get_id", handler.GetIdByMatrix)
	http.HandleFunc("/change_storage", handler.ChangeStorage)
	go func() {
		port := strconv.Itoa(int(config.ServerPort))
		logrus.Infof("Price Management Service is listening on port %s...\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	select {}
}
