package main

import (
	"connector"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"net/http"
	"os"
	"strconv"
	"time"
	"trees"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger    *logrus.Logger
	connector connector.Connector
}

func NewHandler() *Handler {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel) // Set log level to info

	return &Handler{
		logger: logger,
		connector: connector.NewPriceManagerConnector(
			config.PriceManagementHost, strconv.Itoa(int(config.PriceManagementPort)),
			config.RedisHost, config.RedisPassword, config.RedisDB,
		),
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
	response, err := retriever.Search(&info)
	if err != nil {
		h.logServerError(r, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(response)
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
	"Path to the retrieval file .yaml file which contains server_port, price_management_host, "+
		"price_management_port, redis_host, redis_password, redis_db, locations_tree, category_tree,"+
		" segments, db_name_path(Path to data base table and their IDs map). "+
		"Location tree, category tree, segments, db_name_path should be json file")
var config Config
var NoConfig = errors.New("you should set config file. Use --help to information")

const LRUCacheSize = 10000

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

	_, err = trees.BuildLocationTreeFromFile(config.LocationTree)
	if err != nil {
		logrus.Fatal(err)
	}

	_, err = trees.BuildCategoryTreeFromFile(config.CategoryTree)
	if err != nil {
		logrus.Fatal(err)
	}

	err = connector.LoadTableNameByID(config.DBNamePath, true)
	if err != nil {
		logrus.Fatal(err)
	}

	err = connector.LoadTableNameByID(config.BaseTablePath, false)
	if err != nil {
		logrus.Fatal(err)
	}

	err = connector.LoadSegmentsByUserMap(config.Segments)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Config load successfully")

	LRUCache, err = lru.New2Q[CacheKey, CacheValue](LRUCacheSize)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Cache initialized successfully")

	handler := NewHandler()

	http.HandleFunc("/retrieve", handler.PriceRetrievalService)

	go func() {
		port := strconv.Itoa(int(config.ServerPort))
		fmt.Printf("Price Retrieval Service is listening on port %s...\n", port)
		logrus.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	select {}
}
