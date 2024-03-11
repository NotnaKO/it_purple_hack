package main

import (
	"errors"
	"net/http"
	"strconv"
)

type HttpGetRequestInfo struct {
	LocationID, MicrocategoryID, DataBaseID uint64
}

func NewGetRequest(r *http.Request) (HttpGetRequestInfo, error) {
	locationIDStr := r.URL.Query().Get("location_id")
	if locationIDStr == "" {
		return HttpGetRequestInfo{}, errors.New("location_id is required")
	}
	locationID, err := strconv.ParseUint(locationIDStr, 10, 64)
	if err != nil {
		return HttpGetRequestInfo{}, errors.New("invalid location_id")
	}

	microcategoryIDStr := r.URL.Query().Get("microcategory_id")
	if microcategoryIDStr == "" {
		return HttpGetRequestInfo{}, errors.New("microcategory_id is required")
	}
	microcategoryID, err := strconv.ParseUint(microcategoryIDStr, 10, 64)
	if err != nil {
		return HttpGetRequestInfo{}, errors.New("invalid microcategory_id")
	}

	dbIDS := r.URL.Query().Get("data_base_id")
	if dbIDS == "" {
		return HttpGetRequestInfo{}, errors.New("data_base_id is required")
	}
	dbID, err := strconv.ParseUint(dbIDS, 10, 64)
	if err != nil {
		return HttpGetRequestInfo{}, errors.New("invalid microcategory_id")
	}

	return HttpGetRequestInfo{
		LocationID:      locationID,
		MicrocategoryID: microcategoryID,
		DataBaseID:      dbID,
	}, nil
}

type HttpSetRequestInfo struct {
	LocationID, MicrocategoryID uint64
	Price                       float64
	DataBaseID                  uint64
}

func NewSetRequest(r *http.Request) (HttpSetRequestInfo, error) {
	locationIDStr := r.URL.Query().Get("location_id")
	if locationIDStr == "" {
		return HttpSetRequestInfo{}, errors.New("location_id is required")
	}
	locationID, err := strconv.ParseUint(locationIDStr, 10, 64)
	if err != nil {
		return HttpSetRequestInfo{}, errors.New("invalid location_id")
	}

	microcategoryIDStr := r.URL.Query().Get("microcategory_id")
	if microcategoryIDStr == "" {
		return HttpSetRequestInfo{}, errors.New("microcategory_id is required")
	}
	microcategoryID, err := strconv.ParseUint(microcategoryIDStr, 10, 64)
	if err != nil {
		return HttpSetRequestInfo{}, errors.New("invalid microcategory_id")
	}

	priceStr := r.URL.Query().Get("price")
	if priceStr == "" {
		return HttpSetRequestInfo{}, errors.New("price is required")
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return HttpSetRequestInfo{}, errors.New("invalid microcategory_id")
	}

	dbIDS := r.URL.Query().Get("data_base_id")
	if dbIDS == "" {
		return HttpSetRequestInfo{}, errors.New("data_base_id is required")
	}
	dbID, err := strconv.ParseUint(dbIDS, 10, 64)
	if err != nil {
		return HttpSetRequestInfo{}, errors.New("invalid microcategory_id")
	}

	return HttpSetRequestInfo{
		LocationID:      locationID,
		MicrocategoryID: microcategoryID,
		Price:           price,
		DataBaseID:      dbID,
	}, nil
}
