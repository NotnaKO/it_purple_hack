package connector

import "errors"

type SegmentAndTable struct {
	Segment   uint64 // == Table ID
	TableName string
}

type Connector interface {
	GetPrice(locationID, microcategoryID, tableID uint64) (uint64, error)

	GetTablesInOrder(userID uint64) ([]SegmentAndTable, error)
}

var NoResult = errors.New("no answer here")
