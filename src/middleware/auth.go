package middleware

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/util"
)

// WithAuth - authentication middleware
func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authStr := c.Request().Header.Get("Authorization")
		if !strings.Contains(authStr, "Bearer ") {
			return util.JsonError(c, 401, "invalid authentication method")
		}

		token := authStr[7:]
		if token == "" {
			return util.JsonError(c, 401, "missing access_token")
		}

		tokenInfo, err := server.Get().Auth().VerifyAccessToken(token)
		if err != nil {
			return util.JsonError(c, 401, err.Error())
		}

		c.Set("uid", tokenInfo.User.ID)
		return next(c)
	}
}
