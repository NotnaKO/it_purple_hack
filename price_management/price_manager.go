package main

import (
	"database/sql"
)

// TODO create hot table
type PriceManager struct {
	db *sql.DB
}

func NewPriceManagementService(db *sql.DB) *PriceManager {
	return &PriceManager{
		db: db,
	}
}

// SetPrice устанавливает цену для указанных местоположения и микрокатегории

// в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) SetPrice(request *HttpSetRequestInfo) error {
	_, err := p.db.Exec("INSERT INTO price_matrix(location_id, microcategory_id, price) VALUES($1, $2, $3)",
		request.LocationID, request.MicrocategoryID, request.Price)
	return err
}

// GetPrice возвращает цену для указанных местоположения и микрокатегории
// в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) GetPrice(request *HttpGetRequestInfo) (uint64, error) {
	var price uint64
	err := p.db.QueryRow("SELECT price FROM price_matrix WHERE location_id=$1 AND microcategory_id=$2",
		request.LocationID, request.MicrocategoryID).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
