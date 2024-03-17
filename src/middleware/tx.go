package middleware

import (
	"log"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/server"
)

// AddTransaction - middleware that help add transaction to handler
func AddTransaction(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tx := server.Get().DB().NewTransaction()
		c.Set("tx", tx.Tx())

		if rec := recover(); rec != nil {
			log.Printf("rollback transaction, recover: %v\n", rec)
			tx.Rollback()
			return next(c)
		}

		if c.Response().Status == 500 {
			log.Println("rollback transaction, request failed")
			tx.Rollback()
			return next(c)
		}

		if c.Request().Method == "GET" {
			tx.Rollback()
			return next(c)
		}

		tx.Commit()
		return next(c)
	}
}
