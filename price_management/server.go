package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
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

	get_request, err := NewGetRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, err := h.priceManager.GetPrice(&get_request)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Price float64 `json:"price"`
	}{
		Price: price,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	h.logRequest(r, startTime)
}

func (h *Handler) SetPrice(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	set_request, err := NewSetRequest(r)
	if err != nil {
		h.logRequestError(r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.priceManager.SetPrice(&set_request)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

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
	if len(os.Args) != 5 {
		fmt.Println("Usage: ./main user password host dbname")
		return nil, errors.New("Invalid usage")
	}

	user := os.Args[1]
	password := os.Args[2]
	host := os.Args[3]
	dbname := os.Args[4]
	return sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel) // Set log level to info

	db, err := ConnectToDatabase()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	priceManager := NewPriceManagementService(db)
	handler := NewHandler(priceManager, logger)

	http.HandleFunc("/get_price", handler.GetPrice)
	http.HandleFunc("/set_price", handler.SetPrice)

	logger.Info("Price Management Service is running...")
	http.ListenAndServe(":8080", nil)
}
