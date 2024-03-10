package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type PriceManagerConnector struct {
	pmHost string
	pmPort string
}

func NewPriceManagerConnector(pmHost, pmPort string) *PriceManagerConnector {
	return &PriceManagerConnector{
		pmHost: pmHost,
		pmPort: pmPort,
	}
}

func (c *PriceManagerConnector) GetPrice(locationID, microcategoryID uint64) (uint64, error) {
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

	return data.Price, nil
}
