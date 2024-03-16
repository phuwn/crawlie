package db

import (
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

const txKey int = iota

// GetTxFromCtx --
func GetTxFromCtx(c echo.Context) (*gorm.DB, error) {
	if db == nil {
		return nil, ErrNoValidConnection
	}

	tx := c.Get(strconv.Itoa(txKey))
	if tx == nil {
		return nil, ErrTransactionKeyNotValid
	}

	return c.Get(strconv.Itoa(txKey)).(*gorm.DB), nil
}

// SetTxToCtx --
func SetTxToCtx(c echo.Context, tx *gorm.DB) {
	c.Set(strconv.Itoa(txKey), tx)
}
