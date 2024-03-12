package main

import (
	"connector"
	"errors"
	"github.com/hashicorp/golang-lru/v2"
	"github.com/sirupsen/logrus"
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
	Location        uint64 `json:"location"`
	Category        uint64 `json:"category"`
	Price           uint64 `json:"price"`
	DiscountSegment uint64 `json:"tableID"`
	UserID          uint64 `json:"user_id"`
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
	logrus.Debug("Get order: ", segmentWithTable)
	if err != nil {
		return SearchResponse{}, err
	}
	for _, segmentAndTable := range segmentWithTable {
		response, err := r.search(request, request, segmentAndTable.Segment)
		if errors.Is(err, NoSuchCategoryAndLocation) {
			logrus.Debug("No result in this table: continue")
			continue
		}
		if err != nil {
			logrus.Error("Error in search:", err)
			return SearchResponse{}, err
		}
		response.UserID = info.UserID
		response.DiscountSegment = segmentAndTable.Segment
		logrus.Debugf("Return success response: %+v", response)
		return response, nil
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

var LRUCache *lru.TwoQueueCache[CacheKey, CacheValue]

type CacheKey struct {
	searchRequest
	tableID uint64
}

type CacheValue struct {
	resp SearchResponse
	err  error
}

func (r *Retriever) search(request searchRequest, firstRequest searchRequest, tableID uint64) (SearchResponse, error) {
	key := CacheKey{
		searchRequest: request,
		tableID:       tableID,
	}
	if value, ok := LRUCache.Get(key); ok {
		logrus.Debugf("Answer found in LRU cache: %+v", value)
		if value.err != nil {
			return SearchResponse{}, value.err
		} else {
			return value.resp, nil
		}
	}
	logrus.Debug("Answer not found in LRU cache, go to search implementation")
	response, err := r.searchImpl(request, firstRequest, tableID)
	LRUCache.Add(key, CacheValue{
		resp: response,
		err:  err,
	})
	return response, err
}

func (r *Retriever) searchImpl(request searchRequest, firstRequest searchRequest, tableID uint64) (SearchResponse, error) {
	price, err := r.connector.GetPrice(request.location.ID, request.category.ID, tableID)

	if errors.Is(err, connector.NoResult) {
		logrus.Debug("No result from connector")
		nextRequest, err := next(request, firstRequest)
		if err != nil {
			return SearchResponse{}, err
		}
		return r.search(nextRequest, firstRequest, tableID) // TODO: no recursion
	}
	if err != nil {
		logrus.Error("Error in connector:", err)
		return SearchResponse{}, err
	}
	logrus.Debug("Get price from connector:", price)
	response := SearchResponse{
		Location: request.location.ID,
		Category: request.category.ID,
		Price:    price,
	}
	return response, nil
}
