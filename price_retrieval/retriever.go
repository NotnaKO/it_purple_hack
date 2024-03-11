package main

import (
	"connector"
	"errors"
	"trees"
)

type Retriever struct {
	connector connector.Connector
}

type searchRequest struct {
	location *trees.LocationNode
	category *trees.CategoryNode
	userID   uint64
}

type SearchResponse struct {
	location        *trees.LocationNode
	category        *trees.CategoryNode
	price           uint64
	discountSegment uint64
	userID          uint64
}

// Search Возвращает цену в копейках
func (r *Retriever) Search(info *ConnectionInfo) (uint64, error) {
	location := trees.IDToLocationNodeMap[info.LocationID]
	category := trees.IDToCategoryNodeMap[info.MicrocategoryID]
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
	if errors.Is(err, &connector.NoResultError{}) {
		nextRequest, err := next(request, firstRequest)
		if err != nil {
			return 0, err
		}
		return r.search(nextRequest, firstRequest) // TODO: no recursion
	}
	if err != nil {
		return 0, err
	}
	return price, nil
}
