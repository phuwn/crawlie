package config

import (
	"encoding/json"
	"os"
)

var (
	configPath = "./config/config.json"
)

type Router struct {
	ListeningPort int `json:"listening_port,omitempty"`
}

type Authenticator struct {
	JwtSecretKey string `json:"jwt_secret_key,omitempty"`
}

type Database struct {
	Type string `json:"type,omitempty"`
	DSN  string `json:"dsn,omitempty"`
}

type Service struct {
	GoogleOauth2ClientID     string `json:"google_oauth_2_client_id,omitempty"`
	GoogleOauth2ClientSecret string `json:"google_oauth_2_client_secret,omitempty"`
}

type Config struct {
	Router        `json:"router,omitempty"`
	Authenticator `json:"authenticator,omitempty"`
	Database      `json:"database,omitempty"`
	Service       `json:"service,omitempty"`
}

func Load() (*Config, error) {
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
