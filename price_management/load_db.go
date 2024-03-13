package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func (p *PriceManager) loadDB() error {
	err := os.Chdir(config.DataGenerationDirPath)
	if err != nil {
		return err
	}

	logrus.Debugf("Now directory is %s", func(answer string, err error) string {
		if err != nil {
			logrus.Fatal(err)
		}
		return answer
	}(os.Getwd()))

	err = exec.Command("bash", "create_tables.sh").Run()
	if err != nil {
		return err
	}
	logrus.Debug("Master tables created. Now go to partitions.")

	for _, tableName := range p.DataBaseById {
		for i := uint64(0); i < config.CategoriesCount; i += config.TablePartitionSize {
			fullName := fmt.Sprintf("%s.%s", config.DBSchema, tableName)
			request := fmt.Sprintf("CREATE TABLE %s_%d PARTITION OF %s FOR VALUES FROM (%d) TO (%d);",
				fullName, i/config.TablePartitionSize+1, fullName, i, i+config.TablePartitionSize)
			_, err = p.db.Exec(request)
			if err != nil {
				return err
			}
		}
	}
	logrus.Debug("Partitions set successfully. Now go to insert values")
	err = exec.Command("bash", "insert_values.sh").Run()
	if err != nil {
		return err
	}
	logrus.Debug("Databases set successfully")
	return nil
}
