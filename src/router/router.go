package router

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/phuwn/crawlie/src/handler"
	"github.com/phuwn/crawlie/src/handler/keyword"
	"github.com/phuwn/crawlie/src/handler/user"
	"github.com/phuwn/crawlie/src/middleware"
	"github.com/phuwn/crawlie/src/util"
)

// Router - handling routes for incoming request
func Router() *echo.Echo {
	r := echo.New()
	r.HTTPErrorHandler = errorHandler
	r.Use(cors)
	r.GET("/healthz", handler.Healthz)

	v1 := r.Group("/v1")
	r.Pre(mw.RemoveTrailingSlash())
	{
		v1.Use(middleware.AddTransaction)
		v1.POST("/auth", user.SignIn)
		v1.Use(middleware.WithAuth)
	}
	v1.GET("/users/me", user.Get)

	v1.GET("/keywords", keyword.ListByUser)
	v1.GET("/keywords/:name", keyword.Get)
	v1.POST("/keywords", keyword.UploadFile)

	return r
}

// CorsConfig - CORS middleware
var cors = mw.CORSWithConfig(mw.CORSConfig{
	AllowOrigins:     []string{"*"},
	AllowHeaders:     []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowCredentials: true,
})

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
