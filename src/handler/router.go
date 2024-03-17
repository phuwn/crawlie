package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/phuwn/crawlie/src/db"
	"github.com/phuwn/crawlie/src/util"
)

// Router - handling routes for incoming request
func Router(config *Config) *echo.Echo {
	r := echo.New()
	r.HTTPErrorHandler = errorHandler
	r.GET("/healthz", healthz)

	v1 := r.Group("/v1")
	NewAuthenticator(config.JwtSecretKey)
	r.Pre(mw.RemoveTrailingSlash())
	{
		v1.Use(CorsConfig())
		v1.Use(AddTransaction)
		v1.POST("/auth", signIn)
		v1.Use(WithAuth)
	}
	{
		userRoutes(v1)
	}

	return r
}

func healthz(c echo.Context) error {
	err := db.Healthz()
	if err != nil {
		return err
	}
	return c.JSONBlob(200, []byte(`{"message":"ok"}`))
}

// CorsConfig - CORS middleware
func CorsConfig() echo.MiddlewareFunc {
	return mw.CORSWithConfig(mw.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	})
}

func errorHandler(err error, c echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		code := he.Code
		msg := he.Message
		if he.Internal != nil || code == http.StatusInternalServerError {
			log.Printf("%v. %v\n", msg, he.Internal)
			util.JsonError(c, code, "Internal Server Error")
			return
		}
		c.JSON(code, msg)
		return
	}
	log.Println(err.Error())
	util.JsonError(c, http.StatusInternalServerError, "Internal Server Error")
}
