package connector

import (
	"cmp"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"slices"
	"time"

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
	segments := GetSegmentsByUserIDs([]uint64{userID})[userID] // todo: add batching
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

func (c *PriceManagerConnector) GetPrice(locationID, microcategoryID uint64, tableName string) (uint64, error) {
	ctx := context.Background()
	price, err := c.getPriceFromCache(ctx, locationID, microcategoryID, tableName)
	if err != nil {
		// Attempt to fetch price from the Price Manager if it wasn't found in the cache
		price, err = c.fetchPriceFromManager(ctx, locationID, microcategoryID, tableName)
		if err != nil {
			logrus.Error("Error fetching price from Price Manager:", err)
			return 0, err
		}
	}
	return price, nil
}

func (c *PriceManagerConnector) fetchPriceFromManager(ctx context.Context, locationID,
	microcategoryID uint64, tableName string) (uint64, error) {
	url := fmt.Sprintf(
		"http://%s:%s/get_price?location_id=%d&microcategory_id=%d&table_name=%s",
		c.pmHost, c.pmPort, locationID, microcategoryID, tableName)

	resp, err := http.Get(url)
	if errors.Is(err, sql.ErrNoRows) {
		err = errors.Join(err, NoResult)
	}
	if err != nil {
		return 0, err
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
	err = c.setPriceInCache(ctx, locationID, microcategoryID, tableName, data.Price)
	if err != nil {
		logrus.Error("Error setting price in cache:", err)
	}
	return data.Price, nil
}

func (c *PriceManagerConnector) getPriceFromCache(ctx context.Context, locationID, microcategoryID uint64, tableID string) (uint64, error) {
	price, err := c.redisClient.Get(ctx, fmt.Sprintf("%d:%d:%s", locationID, microcategoryID, tableID)).Uint64()
	if errors.Is(err, redis.Nil) {
		// Cache miss
		return 0, errors.New("price not found in cache")
	} else if err != nil {
		return 0, err
	}
	// Cache hit
	return price, nil
}

func (c *PriceManagerConnector) setPriceInCache(ctx context.Context, locationID, microcategoryID uint64, tableName string, price uint64) error {
	return c.redisClient.Set(ctx, fmt.Sprintf("%d:%d:%s", locationID, microcategoryID, tableName), price, time.Hour).Err()
}
