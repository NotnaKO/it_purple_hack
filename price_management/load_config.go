package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServerPort     uint   `yaml:"server_port"`
	PostgresqlUser string `yaml:"postgresql_user"`
	Password       string `yaml:"password"`
	PostgresqlHost string `yaml:"postgresql_host"`
	Dbname         string `yaml:"dbname"`
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
