package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type JSONCategory struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

var DataBaseById map[uint64]string

type Handler struct {
	logger       *logrus.Logger
	priceManager *PriceManager
}

func NewHandler(priceManager *PriceManager, logger *logrus.Logger) *Handler {
	return &Handler{
		logger:          logger,
		priceManager:    priceManager,
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
	user := os.Args[2]
	password := os.Args[3]
	host := os.Args[4]
	dbname := os.Args[5]
	return sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbname))
}

func ParseTableIdJson(filename string)  error {
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
	if err := json.NewDecoder(file).Decode(&tables); err != nil {
		return err
	}
	DataBaseById = make(map[uint64]string)
	for _, table := range tables {
		num, err := strconv.ParseUint(table.Id, 10, 64)
		if err != nil {
			return err 
		}
		DataBaseById[num] = table.Name
	}
	return nil
}

func main() {
	if len(os.Args) != 7 {
		fmt.Println("Usage: ./price_management [server_port] [postgresql_user] [password] [postgresql_host] [dbname] [data_base_json]")
		return
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
	ParseTableIdJson(os.Args[6])
	http.HandleFunc("/get_price", handler.GetPrice)
	http.HandleFunc("/set_price", handler.SetPrice)
	// TODO kubernetes. right now leave only one port
	go func() {
		port := os.Args[1]
		fmt.Printf("Price Management Service is listening on port %s...\n", port)
		http.ListenAndServe(":"+port, nil)
	}()

	select {}
}

/*
go build
./price_management 5432 postgres 1234 localhost price_management data/data_base.json
*/
