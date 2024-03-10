package main

import (
	"database/sql"
	"fmt"
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

// в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) SetPrice(request *HttpSetRequestInfo) error {
	// for debug table id reguest
	//  fmt.Printf("SELECT price FROM %s WHERE location_id=$1 AND microcategory_id=$2", p.DataBaseById[request.DataBaseID])
	_, err := p.db.Exec(fmt.Sprintf("INSERT INTO %s(location_id, microcategory_id, price) VALUES($1, $2, $3)",
						p.DataBaseById[request.DataBaseID]),
		request.LocationID, request.MicrocategoryID, request.Price)
	return err
}

// GetPrice возвращает цену для указанных местоположения и микрокатегории
// в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) GetPrice(request *HttpGetRequestInfo) (uint64, error) {
	var price uint64
	// for debug table id
	// fmt.Printf("SELECT price FROM %s WHERE location_id=$1 AND microcategory_id=$2", p.DataBaseById[request.DataBaseId])
	err := p.db.QueryRow(fmt.Sprintf("SELECT price FROM %s", p.DataBaseById[request.DataBaseID])+"WHERE location_id=$1 AND microcategory_id=$2",
		request.LocationID, request.MicrocategoryID).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
