package handler

import (
	"fmt"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/phuwn/crawlie/src/db"
)

// Router - handling routes for incoming request
func Router(config *Config) *echo.Echo {
	r := echo.New()
	r.GET("/healthz", healthz)

	v1 := r.Group("/v1")
	authenticator := NewAuthenticator(config.JwtSecretKey)
	r.Pre(mw.RemoveTrailingSlash())
	{
		v1.Use(CorsConfig())
		v1.Use(AddTransaction)
		v1.Use(authenticator.WithAuth)
	}
	{
		// userRoutes(r)
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

func jsonError(c echo.Context, code int, msg string) error {
	return c.JSONBlob(code, []byte(fmt.Sprintf(`{"error":"%s"}`, msg)))
}
