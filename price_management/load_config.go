package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	serverPort     uint   `yaml:"server_port"`
	postgresqlUser string `yaml:"postgresql_user"`
	password       string `yaml:"password"`
	postgresqlHost string `yaml:"postgresql_host"`
	dbname         string `yaml:"dbname"`
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
