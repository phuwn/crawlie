package util

import (
	"fmt"

	"github.com/labstack/echo"
)

func JsonError(c echo.Context, code int, msg string) error {
	return c.JSONBlob(code, []byte(fmt.Sprintf(`{"error":"%s"}`, msg)))
}
