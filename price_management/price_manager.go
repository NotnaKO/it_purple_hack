package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

// TODO create hot table
type PriceManager struct {
	db           *sql.DB
	DataBaseById map[uint64]string
}

func NewPriceManagementService(db *sql.DB, filename string) (*PriceManager, error) {
	pm := &PriceManager{
		db: db,
	}
	var err error
	pm.DataBaseById, err = ParseTableIdJson(filename)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

// SetPrice устанавливает цену для указанных местоположения и микрокатегории

// SetPrice в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) SetPrice(request *HttpSetRequestInfo) error {
	// for debug table id reguest
	//  fmt.Printf("SELECT price FROM %s WHERE location_id=$1 AND microcategory_id=$2", p.DataBaseById[request.DataBaseID])
	val, ok := p.DataBaseById[request.DataBaseID]
	if !ok {
		return errors.New("no exist table with that data_base_id")
	}
	_, err := p.db.Exec(fmt.Sprintf("INSERT INTO %s(location_id, microcategory_id, price) VALUES($1, $2, $3)",
		val),
		request.LocationID, request.MicrocategoryID, request.Price)
	return err
}

// GetPrice возвращает цену для указанных местоположения и микрокатегории
// в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) GetPrice(request *HttpGetRequestInfo) (uint64, error) {
	var price uint64
	tableName, ok := p.DataBaseById[request.DataBaseID]
	if !ok {
		return 0, errors.New("no exist table with that data_base_id")
	}
	logrus.Debug(fmt.Sprintf("SELECT price FROM %s WHERE location_id=%d AND microcategory_id=%d\n",
		fmt.Sprintf("%s.%s", config.DBSchema, tableName), request.LocationID, request.MicrocategoryID))
	err := p.db.QueryRow(fmt.Sprintf("SELECT price FROM %s ",
		fmt.Sprintf("%s.%s", config.DBSchema, tableName))+
		"WHERE location_id=$1 AND microcategory_id=$2",
		request.LocationID, request.MicrocategoryID).Scan(&price)
	if err != nil {
		return 0, err
	}
	logrus.Debug("Get price:", price)
	return price, nil
}
