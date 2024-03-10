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
	"github.com/sirupsen/logrus"
)

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
	startTime := time.Now()

	set_request, err := NewSetRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.priceManager.SetPrice(&set_request)
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
	return sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.postgresqlUser, config.password, config.postgresqlHost, config.dbname))
}

var configPath = flag.String("config_path", "",
	"Path to the retrieval file .yaml file which contains server_port,  "+
		"price_management_port, locations_tree, category_tree. Location tree and Category tree should be json file")
var config Config
var NoConfig = errors.New("you should set config file. Use --help to information")

func main() {
	flag.Parse()
	if *configPath == "" {
		logrus.Fatal(NoConfig)
	}
	err := error(nil)
	config, err = loadConfig(*configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel) // Set log level to info

	db, err := ConnectToDatabase()
	if err != nil {
		logger.Fatal(err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Fatal(err)
		}
	}(db)

	priceManager := NewPriceManagementService(db)
	handler := NewHandler(priceManager, logger)

	http.HandleFunc("/get_price", handler.GetPrice)
	http.HandleFunc("/set_price", handler.SetPrice)

	// TODO kubernetes. right now leave only one port
	go func() {
		port := strconv.Itoa(int(config.serverPort))
		fmt.Printf("Price Management Service is listening on port %s...\n", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			logger.Fatal(err)
			return
		}
	}()

	select {}
}
