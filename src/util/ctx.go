package util

import (
	"fmt"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func GetTxFromCtx(c echo.Context) *gorm.DB {
	return c.Get("tx").(*gorm.DB)
}

func GetUserIDFromCtx(c echo.Context) string {
	return fmt.Sprintf("%v", c.Get("uid"))
}
