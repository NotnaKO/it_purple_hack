package price_retrival

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
	logger *logrus.Logger
}

func NewHandler() *Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel) // Set log level to info

	return &Handler{
		logger: logger,
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
	var retriever Retriever
	price, err := retriever.Search(&info)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

func ServerMain() {
	handler := NewHandler()

	http.HandleFunc("/retrieve", handler.PriceRetrievalService)

	// Запускаем HTTP сервер на порту 8080
	fmt.Println("Price Retrieval Service is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
