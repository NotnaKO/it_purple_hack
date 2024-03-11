package connector

import (
	"encoding/json"
	"os"
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
	segmentByUserID = make(map[uint64][]uint64, len(data))
	for _, item := range data {
		segmentByUserID[item.UserID] = item.Segments
	}
	return nil
}

var segmentByUserID map[uint64][]uint64

func GetSegmentsByUserIDs(userIDs []uint64) map[uint64][]uint64 {
	result := make(map[uint64][]uint64, len(userIDs))

	for _, userID := range userIDs {
		user_id_result := segmentByUserID[userID]
		for i := range baseTableName {
			user_id_result = append(user_id_result, i)
		}
		result[userID] = user_id_result
	}
	return result
}

var tableNameByID map[uint64]string
var baseTableName map[uint64]string

type tableAndID struct {
	TableName string `json:"name"`
	ID        uint64 `json:"id"`
}

func LoadTableNameByID(path string, fl bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	var data []tableAndID
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}
	if fl {
		tableNameByID = make(map[uint64]string, len(data))
		for _, item := range data {
			tableNameByID[item.ID] = item.TableName
		}
	} else {
		baseTableName = make(map[uint64]string, len(data))
		for _, item := range data {
			baseTableName[item.ID] = item.TableName
		}
	}
	return nil
}