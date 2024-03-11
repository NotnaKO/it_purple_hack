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
	segmentsByUserID = make(map[uint64][]uint64, len(data))
	for _, item := range data {
		segmentsByUserID[item.UserID] = item.Segments
	}
	return nil
}

var segmentsByUserID map[uint64][]uint64

func GetSegmentsByUserIDs(userIDs []uint64) map[uint64][]uint64 {
	result := make(map[uint64][]uint64, len(userIDs))

	for _, userID := range userIDs {
		userIdResult := segmentsByUserID[userID]
		for i := range baseTableName {
			userIdResult = append(userIdResult, i)
		}
		result[userID] = userIdResult
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
