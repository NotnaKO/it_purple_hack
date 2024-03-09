package price_retrival

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Retriever struct {
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
	resp, err := http.Get(
		fmt.Sprintf("http://localhost:8080/get_price?location_id=%d&microcategory_id=%d",
			request.location.ID, request.category.ID))
	if errors.Is(err, sql.ErrNoRows) {
		nextRequest, err := next(request, firstRequest)
		if err != nil {
			return 0, err
		}
		search, err := r.search(nextRequest, firstRequest)
		if err != nil {
			return 0, err
		}
		return search, nil
	}
	if err != nil {
		return 0, err
	}
	data := resp.Header.Get("Content-Type")
	var handler priceHandler
	err = json.NewDecoder(strings.NewReader(data)).Decode(&handler)
	if err != nil {
		return 0, err
	}
	return handler.price, nil
}
