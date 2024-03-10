package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger    *logrus.Logger
	connector Connector
}

func NewHandler() *Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel) // Set log level to info

	return &Handler{
		logger:    logger,
		connector: NewPriceManagerConnector(strconv.Itoa(int(config.priceManagementHost)), strconv.Itoa(int(config.priceManagementPort))),
	}
}

// PriceRetrievalService обрабатывает запросы retrieve и использует алгоритм RoadUpSearch
func (h *Handler) PriceRetrievalService(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	info, err := NewConnectionInfo(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем функцию roadUpSearch для получения цены с помощью алгоритма RoadUpSearch
	retriever := Retriever{h.connector}
	price, err := retriever.Search(&info)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем ответ в формате JSON
	response := struct {
		Price float64 `json:"price"`
	}{
		Price: float64(price) / 100,
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

var configPath = flag.String("config_path", "",
	"Path to the retrieval file .yaml file which contains server port, price_management_host, "+
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

	_, err = BuildLocationTreeFromFile(config.locationTree)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = BuildCategoryTreeFromFile(config.categoryTree)
	if err != nil {
		logrus.Fatal(err)
	}

	handler := NewHandler()

	http.HandleFunc("/retrieve", handler.PriceRetrievalService)

	go func() {
		port := strconv.Itoa(int(config.serverPort))
		fmt.Printf("Price Retrieval Service is listening on port %s...\n", port)
		logrus.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	select {}
}
