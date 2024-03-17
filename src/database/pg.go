package database

import (
	"errors"

	"github.com/phuwn/crawlie/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNoValidConnection      = errors.New("no valid connection")
	ErrTransactionKeyNotValid = errors.New("transaction key not valid")
)

func NewPostgresConn(cfg config.Database) (Database, error) {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DSN,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := conn.DB()
	if err != nil {
		return nil, err
	}

	return &pgConn{conn: conn, shutdown: db.Close, ping: db.Ping}, nil
}

type pgTx struct {
	tx *gorm.DB
}

func (pg *pgTx) Tx() *gorm.DB {
	return pg.tx
}

func (pg *pgTx) Commit() {
	pg.tx.Commit()
}

func (pg *pgTx) Rollback() {
	pg.tx.Rollback()
}

type pgConn struct {
	conn     *gorm.DB
	shutdown func() error
	ping     func() error
}

func (pg *pgConn) DB() *gorm.DB {
	return pg.conn
}

func (pg *pgConn) NewTransaction() Transaction {
	return &pgTx{pg.conn.Begin()}
}

func (pg *pgConn) Shutdown() error {
	return pg.shutdown()
}

func (pg *pgConn) Healthz() error {
	return pg.ping()
}
