package handler

import (
	"fmt"
	"strings"
	"time"

	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/util"
)

const (
	uidKey = "uid"
)

var auth *authenticator

// TokenInfo - data model contains user's auth info
type TokenInfo struct {
	jwt.StandardClaims
	model.User
}

type authenticator struct {
	jwtSecretKey string
}

func NewAuthenticator(secretKey string) {
	auth = &authenticator{jwtSecretKey: secretKey}
}

// generateJWTToken - create an access_token that represents user's session
func generateJWTToken(info *TokenInfo, expiresAt int64) (string, error) {
	info.ExpiresAt = expiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
	encryptedToken, err := token.SignedString([]byte(auth.jwtSecretKey))
	if err != nil {
		return "", errors.New("failed to sign on token")
	}
	return encryptedToken, nil
}

// verifyAccessToken - validates user's access_token and returns user's id if it's verified
func verifyAccessToken(tokenString string) (*TokenInfo, error) {
	claims := TokenInfo{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.jwtSecretKey), nil
	})

	if !token.Valid {
		return nil, errors.New("token expired, please log out and log in again")
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid signature")
		}
		return nil, errors.New("bad token")
	}
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, errors.New("token expired, please log out and log in again")
	}
	return &claims, nil
}

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

		tokenInfo, err := verifyAccessToken(token)
		if err != nil {
			return util.JsonError(c, 401, err.Error())
		}

		setUserIDToCtx(c, tokenInfo.User.ID)

		return next(c)
	}
}

// setUserIDToCtx - set a uid to echo context for later user's session usage
func setUserIDToCtx(c echo.Context, uid string) {
	c.Set(uidKey, uid)
}

// getUserIDFromCtx - get uid from echo context
func getUserIDFromCtx(c echo.Context) string {
	return fmt.Sprintf("%v", c.Get(uidKey))
}
