package connector

import (
	"cmp"
	"encoding/json"
	"errors"
	"os"
	"slices"
)

type userAndSegments struct {
	UserID   uint64   `json:"user_id"`
	Segments []uint64 `json:"segments"`
}

func LoadSegmentsByUserMap(path string) error {
	if tableNameByID == nil {
		return errors.New("load tableNameByID before SegmentsByUser")
	}
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
	discountTablesByUserID = make(map[uint64][]SegmentAndTable, len(data))
	for _, item := range data {
		currentSegments := make([]SegmentAndTable, len(item.Segments))
		for i, segment := range item.Segments {
			currentSegments[i] = SegmentAndTable{
				Segment:   segment,
				TableName: tableNameByID[segment],
			}
		}
		slices.SortFunc(currentSegments, func(a, b SegmentAndTable) int {
			return -cmp.Compare(a.TableName, b.TableName)
		})
		discountTablesByUserID[item.UserID] = currentSegments
	}
	return nil
}

var discountTablesByUserID map[uint64][]SegmentAndTable

var tableNameByID map[uint64]string
var baselineTables []SegmentAndTable

func LoadTableNameByID(path string, isDiscountLoad bool) error {
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
	if isDiscountLoad {
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
