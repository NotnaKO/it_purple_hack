package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerPort          uint   `yaml:"server_port"`
	PriceManagementHost string `yaml:"price_management_host"`
	PriceManagementPort uint   `yaml:"price_management_port"`
	LocationTree        string `yaml:"location_tree"`
	CategoryTree        string `yaml:"category_tree"`
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
