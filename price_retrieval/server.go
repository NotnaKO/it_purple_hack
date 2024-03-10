package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
		connector: NewPriceManagerConnector(os.Args[2], os.Args[3]),
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
	json.NewEncoder(w).Encode(response)

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

func main() {
	if len(os.Args) != 6 {
		fmt.Println("Usage: ./price_retrieval [server_port] [price_management_host] [price_management_port] [locations_tree.json] [category_tree.json]")
		return
	}

	locationTreeJson := os.Args[4]
	BuildLocationTreeFromFile(locationTreeJson)

	categoryTreeJson := os.Args[5]
	BuildCategoryTreeFromFile(categoryTreeJson)

	handler := NewHandler()

	http.HandleFunc("/retrieve", handler.PriceRetrievalService)

	go func() {
		port := os.Args[1]
		fmt.Printf("Price Retrieval Service is listening on port %s...\n", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	select {}
}
