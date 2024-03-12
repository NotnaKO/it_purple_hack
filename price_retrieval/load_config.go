package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerPort          uint   `yaml:"server_port"`
	PriceManagementHost string `yaml:"price_management_host"`
	PriceManagementPort uint   `yaml:"price_management_port"`
	RedisHost           string `yaml:"redis_host"`
	RedisPassword       string `yaml:"redis_password"`
	RedisDB             int    `yaml:"redis_db"`
	LocationTree        string `yaml:"location_tree"`
	CategoryTree        string `yaml:"category_tree"`
	Segments            string `yaml:"segments"`
	DBNamePath          string `yaml:"db_name_path"`
	BaseTablePath       string `yaml:"base_table"`
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
