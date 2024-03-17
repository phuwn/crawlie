package util

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func JsonError(c echo.Context, code int, msg string) error {
	return c.JSONBlob(code, []byte(fmt.Sprintf(`{"error":"%s"}`, msg)))
}

func GetTxFromCtx(c echo.Context) *gorm.DB {
	return c.Get("tx").(*gorm.DB)
}

func GetUserIDFromCtx(c echo.Context) string {
	return fmt.Sprintf("%v", c.Get("uid"))
}

func JsonParse(c echo.Context, v interface{}) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return JsonError(c, 400, "unable to read the request body")
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return JsonError(c, 400, "wrong request form data")
	}
	return nil
}
