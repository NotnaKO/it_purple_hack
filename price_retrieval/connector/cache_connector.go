package connector

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/sirupsen/logrus"
)

type PriceManagerConnector struct {
	pmHost      string
	pmPort      string
	redisClient *redis.Client
}

/*
GetTablesInOrder
WARNING: this algorithm assume that baseline tables have names as baseline_matrix_Nand discount table have names discount_matrix_N. Other names will give unspecified result
*/
func (c *PriceManagerConnector) GetTablesInOrder(userID uint64) ([]SegmentAndTable, error) {
	// In this problem we mustn't optimize this by caching
	segments := discountTablesByUserID[userID]
	answer := make([]SegmentAndTable, len(segments))
	for i, item := range segments {
		answer[i] = SegmentAndTable{
			Segment:   item,
			TableName: tableNameByID[item],
		}
	}
	slices.SortFunc(answer, func(a, b SegmentAndTable) int {
		return -cmp.Compare(a.TableName, b.TableName)
	})
	answer = append(answer, baselineTables...)
	// Sort by inverse alphabetic order,
	// so discount with the higher than others and discount higher than baseline
	return answer, nil
}

var _ Connector = &PriceManagerConnector{} // interface check

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

func (c *PriceManagerConnector) GetPrice(locationID, microCategoryID, tableID uint64) (uint64, error) {
	ctx := context.Background()
	price, err := c.getPriceFromCache(ctx, locationID, microCategoryID, tableID)
	if err != nil {
		// Attempt to fetch price from the Price Manager if it wasn't found in the cache
		price, err = c.fetchPriceFromManager(ctx, locationID, microCategoryID, tableID)
		if err != nil {
			logrus.Error("Error fetching price from Price Manager:", err)
			return 0, err
		}
	}
	return price, nil
}

func (c *PriceManagerConnector) fetchPriceFromManager(ctx context.Context, locationID,
	microcategoryID, tableID uint64) (uint64, error) {
	logrus.Debug("Request from manager: ", fmt.Sprintf(
		"http://%s:%s/get_price?location_id=%d&microcategory_id=%d&data_base_id=%d",
		c.pmHost, c.pmPort, locationID, microcategoryID, tableID))
	url := fmt.Sprintf(
		"http://%s:%s/get_price?location_id=%d&microcategory_id=%d&data_base_id=%d",
		c.pmHost, c.pmPort, locationID, microcategoryID, tableID)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == http.StatusInternalServerError {
		text, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error("error while reading body:", err)
			return 0, err
		}
		if bytes.Contains(text, []byte("no rows in result set")) {
			logrus.Debug("No rows in table: ", tableID)
			return 0, NoResult
		}
		logrus.Debug("Just internal error")
		return 0, errors.New("request failed with status: " + resp.Status)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error("Error while closing response body:", err)
		}
	}(resp.Body)

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
	err = c.setPriceInCache(ctx, locationID, microcategoryID, tableID, data.Price)
	if err != nil {
		logrus.Error("Error setting price in cache:", err)
	}
	return data.Price, nil
}

func (c *PriceManagerConnector) getPriceFromCache(ctx context.Context,
	locationID, microCategoryID, tableID uint64) (uint64, error) {
	logrus.Debug("Search in cache by: ", fmt.Sprintf("%d:%d:%d",
		locationID, microCategoryID, tableID))
	price, err := c.redisClient.Get(ctx, fmt.Sprintf("%d:%d:%d",
		locationID, microCategoryID, tableID)).Uint64()
	if errors.Is(err, redis.Nil) {
		// Cache miss
		logrus.Debug(fmt.Sprintf("%d:%d:%d",
			locationID, microCategoryID, tableID), " not found in cache")
		return 0, errors.New("price not found in cache")
	} else if err != nil {
		logrus.Error("Cache finding", fmt.Sprintf("%d:%d:%d",
			locationID, microCategoryID, tableID), "error:", err)
		return 0, err
	}
	logrus.Debug("find in cache:", price)
	// Cache hit
	return price, nil
}

func (c *PriceManagerConnector) setPriceInCache(ctx context.Context,
	locationID, microCategoryID, tableID uint64, price uint64) error {
	return c.redisClient.Set(ctx, fmt.Sprintf("%d:%d:%d",
		locationID, microCategoryID, tableID), price, time.Hour).Err()
}
