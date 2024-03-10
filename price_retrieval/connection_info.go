package price_retrival

import (
	"errors"
	"net/http"
	"strconv"
)

type ConnectionInfo struct {
	LocationID, MicrocategoryID, UserID uint64
}

func NewConnectionInfo(r *http.Request) (ConnectionInfo, error) {
	locationIDStr := r.URL.Query().Get("location_id")
	if locationIDStr == "" {
		return ConnectionInfo{}, errors.New("location_id is required")
	}
	locationID, err := strconv.ParseUint(locationIDStr, 10, 64)
	if err != nil {
		return ConnectionInfo{}, errors.New("invalid location_id")
	}

	microcategoryIDStr := r.URL.Query().Get("microcategory_id")
	if microcategoryIDStr == "" {
		return ConnectionInfo{}, errors.New("microcategory_id is required")
	}
	microcategoryID, err := strconv.ParseUint(microcategoryIDStr, 10, 64)
	if err != nil {
		return ConnectionInfo{}, errors.New("invalid microcategory_id")
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		return ConnectionInfo{}, errors.New("user_id is required")
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return ConnectionInfo{}, errors.New("invalid user_id")
	}

	return ConnectionInfo{
		LocationID:      locationID,
		MicrocategoryID: microcategoryID,
		UserID:          userID,
	}, nil
}
