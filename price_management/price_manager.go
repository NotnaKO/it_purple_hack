package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

type PriceManager struct {
	db           *sql.DB
	DataBaseById map[uint64]string
	mx           sync.Mutex
	DataModify   map[uint64]bool
}

func NewPriceManagementService(db *sql.DB, filename string) (*PriceManager, error) {
	pm := &PriceManager{
		db: db, DataModify: make(map[uint64]bool),
	}
	err := pm.ParseTableIdJson(filename)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

func (p *PriceManager) ChangeStorage(request *HttpChangeStorage) (bool, error) {
	tableName := request.DataBaseName
	find_matrx := false
	matrix_id := uint64(0)
	mx_id := uint64(0)
	for i, j := range p.DataBaseById {
		mx_id = max(mx_id, i)
		if j == tableName {
			matrix_id = i
			find_matrx = true
			break
		}
	}
	p.mx.Lock()
	if !find_matrx {
		matrix_id = mx_id + 1
		p.DataBaseById[matrix_id] = p.DataBaseById[1]
		p.DataBaseById[1] = tableName
		p.createTable(tableName)
		
	} else {
		p.DataBaseById[matrix_id], p.DataBaseById[1] = p.DataBaseById[1], p.DataBaseById[matrix_id]
	}
	p.dumpTables()
	p.mx.Unlock()
	return (tableName != ""), nil
}

// SetPrice устанавливает цену для указанных местоположения и микрокатегории

var SetError = errors.New("error in set in saving data")

// SetPrice в поле DataBaseId возвращается id таблицы из json таблиц
func (p *PriceManager) SetPrice(request *HttpSetRequestInfo) error {
	tableName, ok := p.DataBaseById[request.DataBaseID]
	if !ok {
		return errors.New("no exist table with that data_base_id")
	}
	tableNameInRequest := fmt.Sprintf("%s.%s", config.DBSchema, tableName)

	// Flush cache
	err := exec.Command("redis-cli", "flushall").Run()
	if err != nil {
		logrus.Error("Error while flush cache:", err)
		return errors.Join(SetError, err)
	} else {
		logrus.Debug("Cache flushed successfully")
	}

	requestToDB := fmt.Sprintf(
		"INSERT INTO %s(location_id, microcategory_id, price) VALUES(%d, %d, %d) ON CONFLICT ON CONSTRAINT pk_%s DO UPDATE SET price=%d\n",
		tableNameInRequest,
		request.LocationID, request.MicroCategoryID, request.Price, tableName, request.Price)
	logrus.Debug(requestToDB)
	_, err = p.db.Exec(requestToDB)

	p.DataModify[request.DataBaseID] = true
	logrus.Debugf("Set modification to %s", tableName)
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
		p.mx.Lock()
		matrix_id = mx_id + 1
		p.DataBaseById[matrix_id] = request.DataBaseName
		err := p.createTable(request.DataBaseName)
		if err != nil {
			delete(p.DataBaseById, matrix_id)
			logrus.Debug(err)
			p.mx.Unlock()
			return 0, nil
		}
		p.dumpTables()
		p.mx.Unlock()
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

func (p *PriceManager) waitTimer(tickChannel <-chan time.Time) {
	for {
		<-tickChannel
		err := p.dumpTables()
		if err != nil {
			logrus.Error("Couldn't dump tables by timer, because get error: ", err)
		}
	}
}

func (p *PriceManager) dumpTables() error {
	beginPath, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(config.DataGenerationDirPath)
	if err != nil {
		return err
	}
	for id, modify := range p.DataModify {
		if modify {
			logrus.Debugf("Statrting dumping table %s", p.DataBaseById[id])
			tableName := fmt.Sprintf("%s.%s", config.DBSchema, p.DataBaseById[id])
			fileName := fmt.Sprintf("%s.sql", p.DataBaseById[id])
			for i := uint64(0); i < config.CategoriesCount; i += config.TablePartitionSize {
				fullName := fmt.Sprintf("%s_%d", tableName, i/config.TablePartitionSize+1)
				cmd := exec.Command("bash", "dump.sh",
					fullName, config.Dbname, fileName)
				err := cmd.Run()
				if err != nil {
					return err
				}
			}
			p.DataModify[id] = false
			logrus.Debugf("End dumping table %s", p.DataBaseById[id])

		}
	}
	err = os.Chdir(beginPath)
	if err != nil {
		return err
	}

	return nil
}
