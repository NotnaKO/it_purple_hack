package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func (p *PriceManager) createTable(tableName string) error {
	for i := uint64(0); i < config.CategoriesCount; i += config.TablePartitionSize {
		fullName := fmt.Sprintf("%s.%s", config.DBSchema, tableName)
		request := fmt.Sprintf("CREATE TABLE %s_%d PARTITION OF %s FOR VALUES FROM (%d) TO (%d);",
			fullName, i/config.TablePartitionSize+1, fullName, i, i+config.TablePartitionSize)
		_, err := p.db.Exec(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PriceManager) loadDB() error {
	beginPath, err := os.Getwd()
	if err != nil {
		logrus.Error("Cannot get current path")
		return err
	}
	err = os.Chdir(config.DataGenerationDirPath)
	if err != nil {
		return err
	}

	logrus.Debugf("Now directory is %s", func(answer string, err error) string {
		if err != nil {
			logrus.Fatal(err)
		}
		return answer
	}(os.Getwd()))

	cmd := exec.Command("bash", "create_tables.sh")
	var out strings.Builder
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		logrus.Error(out.String())
		return err
	}
	logrus.Debug("Master tables created. Now go to partitions.")

	for _, tableName := range p.DataBaseById {
		err := p.createTable(tableName)
		if err != nil {
			return err
		}
	}
	logrus.Debug("Partitions set successfully. Now go to insert values")
	err = exec.Command("bash", "insert_values.sh").Run()
	if err != nil {
		return err
	}
	err = os.Chdir(beginPath)
	if err != nil {
		return err
	}
	logrus.Debug("Databases set successfully")
	return nil
}
