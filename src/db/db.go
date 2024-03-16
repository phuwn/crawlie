package db

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db                        *gorm.DB
	ErrNoValidConnection      = errors.New("no valid connection")
	ErrTransactionKeyNotValid = errors.New("transaction key not valid")
)

func NewPostgresConn(config *Config) error {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	db = conn
	return nil
}

func Get() *gorm.DB {
	return db
}

func Healthz() error {
	if db == nil {
		return ErrNoValidConnection
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
