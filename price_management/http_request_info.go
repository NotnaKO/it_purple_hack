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

type HttpGetMatrixByIdRequestInfo struct {
	DataBaseID uint64
}

func NewGetMatrixByIdRequest(r *http.Request) (HttpGetMatrixByIdRequestInfo, error) {

	dbIDS := r.URL.Query().Get("data_base_id")
	if dbIDS == "" {
		return HttpGetMatrixByIdRequestInfo{}, errors.New("data_base_id is required")
	}
	dbID, err := strconv.ParseUint(dbIDS, 10, 64)
	if err != nil {
		return HttpGetMatrixByIdRequestInfo{}, errors.New("invalid data_base_id")
	}

	return HttpGetMatrixByIdRequestInfo{
		DataBaseID: dbID,
	}, nil
}

type HttpGetIdByMatrixRequestInfo struct {
	DataBaseName string
}

func NewGetIdByMatrixRequest(r *http.Request) (HttpGetIdByMatrixRequestInfo, error) {

	dbIDS := r.URL.Query().Get("data_base_name")
	if dbIDS == "" {
		return HttpGetIdByMatrixRequestInfo{}, errors.New("data_base_name is required")
	}

	return HttpGetIdByMatrixRequestInfo{
		DataBaseName: dbIDS,
	}, nil
}

type HttpSetRequestInfo struct {
	LocationID, MicroCategoryID, Price, DataBaseID uint64
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

	microCategoryIDStr := r.URL.Query().Get("microcategory_id")
	if microCategoryIDStr == "" {
		return HttpSetRequestInfo{}, errors.New("microcategory_id is required")
	}
	microCategoryID, err := strconv.ParseUint(microCategoryIDStr, 10, 64)
	if err != nil {
		return HttpSetRequestInfo{}, errors.New("invalid microcategory_id")
	}

	priceStr := r.URL.Query().Get("price")
	if priceStr == "" {
		return HttpSetRequestInfo{}, errors.New("price is required")
	}
	priceInFloat, err := strconv.ParseFloat(priceStr, 64)
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

	price := uint64(priceInFloat * 100) // convert float

	return HttpSetRequestInfo{
		LocationID:      locationID,
		MicroCategoryID: microCategoryID,
		Price:           price,
		DataBaseID:      dbID,
	}, nil
}
