package handler

import (
	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/server"
)

func Healthz(c echo.Context) error {
	err := server.Get().DB().Healthz()
	if err != nil {
		return err
	}
	return c.JSONBlob(200, []byte(`{"message":"ok"}`))
}
