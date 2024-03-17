package database

import (
	"gorm.io/gorm"
)

type Database interface {
	DB() *gorm.DB
	NewTransaction() Transaction
	Shutdown() error
	Healthz() error
}

type Transaction interface {
	Tx() *gorm.DB
	Commit()
	Rollback()
}
