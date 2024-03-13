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
	tableName, ok := p.DataBaseById[request.DataBaseID]
	if !ok {
		return errors.New("no exist table with that data_base_id")
	}
	tableNameInRequest := fmt.Sprintf("%s.%s", config.DBSchema, tableName)
	requestToDB := fmt.Sprintf(
		"INSERT INTO %s(location_id, microcategory_id, price) VALUES(%d, %d, %d) ON CONFLICT ON CONSTRAINT pk_%s DO UPDATE SET price=%d\n",
		tableNameInRequest,
		request.LocationID, request.MicroCategoryID, request.Price, tableName, request.Price)
	logrus.Debug(requestToDB)
	_, err := p.db.Exec(requestToDB)
	return err
}

func (p *PriceManager) GetMatrixById(request *HttpGetMatrixByIdRequestInfo) (string, error) {
	tableName, ok := p.DataBaseById[request.DataBaseID]
	if !ok {
		return "no exist table with that data_base_id", errors.New("no exist table with that data_base_id")
	}
	return tableName, nil
}

func (p *PriceManager) GetIdByMatrix(request *HttpGetIdByMatrixRequestInfo) (uint64, error) {
	find_matrx := false
	matrix_id := uint64(0)
	mx_id := uint64(0)
	for i, j := range p.DataBaseById {
		mx_id = max(mx_id, i)
		if j == request.DataBaseName {
			matrix_id = i
			find_matrx = true
			break
		}
	}
	if !find_matrx {
		matrix_id = mx_id + 1
		p.DataBaseById[matrix_id] = request.DataBaseName
		err := p.loadDB()
		if err != nil {
			delete(p.DataBaseById, matrix_id)
			logrus.Debug(err)
			p.loadDB()
			return 0, nil
		}
	}
	return matrix_id, nil
}

// GetPrice возвращает цену для указанных местоположения и микрокатегории
// в поле DataBaseId возвращается id таблицы из json таблиц(по умолчанию 1)
func (p *PriceManager) GetPrice(request *HttpGetRequestInfo) (uint64, error) {
	var price uint64
	tableName, ok := p.DataBaseById[request.DataBaseID]
	if !ok {
		return 0, errors.New("no exist table with that data_base_id")
	}
	requestToDB := fmt.Sprintf("SELECT price FROM %s_%d WHERE location_id=%d AND microcategory_id=%d\n",
		fmt.Sprintf("%s.%s", config.DBSchema, tableName),
		request.MicrocategoryID/config.TablePartitionSize+1, request.LocationID, request.MicrocategoryID)
	logrus.Debug(requestToDB)
	err := p.db.QueryRow(requestToDB).Scan(&price)
	if err != nil {
		return 0, err
	}
	logrus.Debug("Get price:", price)
	return price, nil
}
