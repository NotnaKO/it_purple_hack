package main

type Connector interface {
	// TODO support multiple tables
	GetPrice(locationID, microcategoryID, discountSegment uint64) (uint64, error)
}

type NoResultError struct{}

func (e *NoResultError) Error() string {
	return "invalid request"
}
