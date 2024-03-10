package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	serverPort          uint   `yaml:"server_port"`
	priceManagementHost uint   `yaml:"price_management_host"`
	priceManagementPort uint   `yaml:"price_management_port"`
	locationTree        string `yaml:"location_tree"`
	categoryTree        string `yaml:"category_tree"`
}

func loadConfig(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
