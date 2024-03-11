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
}

type SearchResponse struct {
	location        *trees.LocationNode
	category        *trees.CategoryNode
	price           uint64
	discountSegment uint64
	userID          uint64
}

var NotFound = errors.New("no prices for this request has found")

// Search Возвращает цену в копейках
func (r *Retriever) Search(info *ConnectionInfo) (SearchResponse, error) {
	location := trees.IDToLocationNodeMap[info.LocationID]
	category := trees.IDToCategoryNodeMap[info.MicrocategoryID]
	if location == nil || category == nil {
		return SearchResponse{}, NoSuchCategoryAndLocation
	}

	request := searchRequest{
		location: location,
		category: category,
	}
	segmentWithTable, err := r.connector.GetTablesInOrder(info.UserID)
	if err != nil {
		return SearchResponse{}, err
	}
	for _, segmentAndTable := range segmentWithTable {
		response, err := r.search(request, request, segmentAndTable.TableName)
		if errors.Is(err, connector.NoResult) {
			continue
		}
		if err != nil {
			return SearchResponse{}, err
		}
		response.userID = info.UserID
		response.discountSegment = segmentAndTable.Segment
	}
	return SearchResponse{}, NotFound
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

func (r *Retriever) search(request searchRequest, firstRequest searchRequest, tableName string) (SearchResponse, error) {
	price, err := r.connector.GetPrice(request.location.ID, request.category.ID, tableName)
	if errors.Is(err, connector.NoResult) {
		nextRequest, err := next(request, firstRequest)
		if err != nil {
			return SearchResponse{}, err
		}
		return r.search(nextRequest, firstRequest, tableName) // TODO: no recursion
	}
	if err != nil {
		return SearchResponse{}, err
	}
	response := SearchResponse{
		location: request.location,
		category: request.category,
		price:    price,
	}
	return response, nil
}
