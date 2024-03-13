package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerPort            uint   `yaml:"server_port"`
	PostgresqlUser        string `yaml:"postgresql_user"`
	Password              string `yaml:"password"`
	PostgresqlHost        string `yaml:"postgresql_host"`
	Dbname                string `yaml:"dbname"`
	DBSchema              string `yaml:"db_schema"`
	DbPathName            string `yaml:"db_path_name"`
	TablePartitionSize    uint64 `yaml:"table_partition_size"`
	DataGenerationDirPath string `yaml:"data_generation_dir_path"`
	CategoriesCount       uint64 `yaml:"categories_count"`
}

func loadConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
