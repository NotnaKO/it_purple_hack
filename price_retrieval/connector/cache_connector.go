package connector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type PriceManagerConnector struct {
	pmHost string
	pmPort string
}

/*
GetTablesInOrder
WARNING: this algorithm assume that baseline tables have names as baseline_matrix_Nand discount table have names discount_matrix_N. Other names will give unspecified result
*/
func (c *PriceManagerConnector) GetTablesInOrder(userID uint64) ([]SegmentAndTable, error) {
	// In this problem we mustn't optimize this by caching
	segments := discountTablesByUserID[userID]
	answer := make([]SegmentAndTable, 0, len(segments)+len(baselineTables))
	answer = append(answer, segments...)
	answer = append(answer, baselineTables...)
	return answer, nil
}

var _ Connector = &PriceManagerConnector{} // interface check

func NewPriceManagerConnector(pmHost, pmPort string) *PriceManagerConnector {
	return &PriceManagerConnector{
		pmHost: pmHost,
		pmPort: pmPort,
	}
}

func (c *PriceManagerConnector) GetPrice(locationID, microCategoryID, tableID uint64) (uint64, error) {
	price, err := c.fetchPriceFromManager(locationID, microCategoryID, tableID)
	if err != nil {
		logrus.Error("Error fetching price from Price Manager:", err)
		return 0, err
	}
	return price, nil
}

func (c *PriceManagerConnector) fetchPriceFromManager(locationID, microCategoryID, tableID uint64) (uint64, error) {
	logrus.Debug("Request from manager: ", fmt.Sprintf(
		"http://%s:%s/get_price?location_id=%d&microcategory_id=%d&data_base_id=%d",
		c.pmHost, c.pmPort, locationID, microCategoryID, tableID))
	url := fmt.Sprintf(
		"http://%s:%s/get_price?location_id=%d&microcategory_id=%d&data_base_id=%d",
		c.pmHost, c.pmPort, locationID, microCategoryID, tableID)

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
	return data.Price, nil
}
