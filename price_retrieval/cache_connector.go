package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type PriceManagerConnector struct {
	pmHost      string
	pmPort      string
	redisClient *redis.Client
}

func NewPriceManagerConnector(pmHost, pmPort, redisAddr, redisPassword string, redisDB int) *PriceManagerConnector {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &PriceManagerConnector{
		pmHost:      pmHost,
		pmPort:      pmPort,
		redisClient: redisClient,
	}
}

func (c *PriceManagerConnector) GetPrice(locationID, microcategoryID uint64) (uint64, error) {
	ctx := context.Background()
	price, err := c.getPriceFromCache(ctx, locationID, microcategoryID)
	if err != nil {
		//logrus.Error("Error retrieving price from cache:", err)

		// Attempt to fetch price from the Price Manager if it wasn't found in the cache
		price, err = c.fetchPriceFromManager(ctx, locationID, microcategoryID)
		if err != nil {
			logrus.Error("Error fetching price from Price Manager:", err)
			return 0, err
		}
	}
	return price, nil
}

func (c *PriceManagerConnector) fetchPriceFromManager(ctx context.Context, locationID, microcategoryID uint64) (uint64, error) {
	url := fmt.Sprintf("http://%s:%s/get_price?location_id=%d&microcategory_id=%d", c.pmHost, c.pmPort, locationID, microcategoryID)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("request failed with status: " + resp.Status)
	}

	var data struct {
		Price uint64 `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	// Store the fetched price in the cache
	err = c.setPriceInCache(ctx, locationID, microcategoryID, data.Price)
	if err != nil {
		logrus.Error("Error setting price in cache:", err)
	}
	return data.Price, nil
}

func (c *PriceManagerConnector) getPriceFromCache(ctx context.Context, locationID, microcategoryID uint64) (uint64, error) {
	price, err := c.redisClient.Get(ctx, fmt.Sprintf("%d:%d", locationID, microcategoryID)).Uint64()
	if err == redis.Nil {
		// Cache miss
		return 0, errors.New("price not found in cache")
	} else if err != nil {
		return 0, err
	}
	// Cache hit
	return price, nil
}

func (c *PriceManagerConnector) setPriceInCache(ctx context.Context, locationID, microcategoryID, price uint64) error {
	return c.redisClient.Set(ctx, fmt.Sprintf("%d:%d", locationID, microcategoryID), price, time.Hour).Err()
}
