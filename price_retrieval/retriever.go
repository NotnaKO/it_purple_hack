package main

import (
	"errors"
)

type Retriever struct {
	connector Connector
}

type searchRequest struct {
	location *LocationNode
	category *CategoryNode
	userID   uint64
}

// Search Возвращает цену в копейках
func (r *Retriever) Search(info *ConnectionInfo) (uint64, error) {
	location := IDToLocationNodeMap[info.LocationID]
	category := IDToCategoryNodeMap[info.MicrocategoryID]
	if location == nil || category == nil {
		return 0, NoSuchCategoryAndLocation
	}

	request := searchRequest{
		location: location,
		category: category,
		userID:   info.UserID,
	}
	return r.search(request, request)
}

type priceHandler struct {
	price uint64
}

var NoSuchCategoryAndLocation = errors.New("no such category and location")

func next(request searchRequest, first searchRequest) (searchRequest, error) {
	if request.category.Parent != nil {
		return searchRequest{location: request.location, category: request.category.Parent}, nil
	}
	if request.location.Parent != nil {
		return searchRequest{location: request.location.Parent, category: first.category}, nil
	}
	return request, NoSuchCategoryAndLocation
}

func (r *Retriever) search(request searchRequest, firstRequest searchRequest) (uint64, error) {
	// TODO add discount tables
	price, err := r.connector.GetPrice(request.location.ID, request.category.ID)
	if errors.Is(err, &NoResultError{}) {
		nextRequest, err := next(request, firstRequest)
		if err != nil {
			return 0, err
		}
		return r.search(nextRequest, firstRequest)
	}
	if err != nil {
		return 0, err
	}
	return price, nil
}
