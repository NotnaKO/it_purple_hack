package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

// TODO create hot table
type PriceManager struct {
	db           *sql.DB
	DataBaseById map[uint64]string
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

var SetError = errors.New("error in set in saving data")

// SetPrice в поле DataBaseId возвращается id таблицы из json таблиц
func (p *PriceManager) SetPrice(request *HttpSetRequestInfo) error {
	// for debug table id reguest
	//  fmt.Printf("SELECT price FROM %s WHERE location_id=$1 AND microcategory_id=$2", p.DataBaseById[request.DataBaseID])
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
