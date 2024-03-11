package connector

import (
	"cmp"
	"encoding/json"
	"os"
	"slices"
)

type userAndSegments struct {
	UserID   uint64   `json:"user_id"`
	Segments []uint64 `json:"segments"`
}

func LoadSegmentsByUserMap(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var data []userAndSegments
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
	discountTablesByUserID = make(map[uint64][]uint64, len(data))
	for _, item := range data {
		discountTablesByUserID[item.UserID] = item.Segments
	}
	return nil
}

var discountTablesByUserID map[uint64][]uint64

var tableNameByID map[uint64]string
var baselineTables []SegmentAndTable

type tableAndID struct {
	TableName string `json:"name"`
	ID        uint64 `json:"id"`
}

func LoadTableNameByID(path string, is_discount_load bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var data []SegmentAndTable
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
	if is_discount_load {
		tableNameByID = make(map[uint64]string, len(data))
		for _, item := range data {
			tableNameByID[item.Segment] = item.TableName
		}
	} else {
		baselineTables = data
		slices.SortFunc(baselineTables, func(a, b SegmentAndTable) int {
			return -cmp.Compare(a.TableName, b.TableName)
		})
	}
	return nil
}
