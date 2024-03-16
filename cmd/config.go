package main

import (
	"encoding/json"
	"os"

	"github.com/phuwn/crawlie/src/db"
	"github.com/phuwn/crawlie/src/handler"
	"github.com/phuwn/crawlie/src/service"
)

var (
	configPath = "./config/config.json"
)

type Config struct {
	Service  *service.Config `json:"service"`
	Server   *handler.Config `json:"server"`
	Database *db.Config      `json:"database"`
}

func configLoad() (*Config, error) {
	setConfigPath := os.Getenv("CONFIG_PATH")
	if setConfigPath != "" {
		configPath = setConfigPath
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	return &config, json.Unmarshal(f, &config)
}
