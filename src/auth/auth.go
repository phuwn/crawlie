package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/phuwn/crawlie/src/config"
	"github.com/phuwn/crawlie/src/model"
)

var (
	ErrTokenExpired      = errors.New("token expired, please log out and log in again")
	ErrInvalidSignature  = errors.New("invalid signature")
	ErrBadToken          = errors.New("bad token")
	ErrFailedToSignToken = errors.New("failed to sign on token")
)

// TokenInfo - data model contains user's auth info
type TokenInfo struct {
	jwt.StandardClaims
	model.User
}

func NewAuthenticator(cfg config.Authenticator) *Authenticator {
	return &Authenticator{cfg.JwtSecretKey}
}

type Authenticator struct {
	jwtSecretKey string
}

// GenerateJWTToken - create an access_token that represents user's session
func (auth *Authenticator) GenerateJWTToken(info *TokenInfo, expiresAt int64) (string, error) {
	info.ExpiresAt = expiresAt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
	encryptedToken, err := token.SignedString([]byte(auth.jwtSecretKey))
	if err != nil {
		return "", ErrFailedToSignToken
	}
	return encryptedToken, nil
}

// VerifyAccessToken - validates user's access_token and returns user's id if it's verified
func (auth *Authenticator) VerifyAccessToken(tokenString string) (*TokenInfo, error) {
	claims := TokenInfo{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.jwtSecretKey), nil
	})

	if !token.Valid {
		return nil, ErrTokenExpired
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, ErrInvalidSignature
		}
		return nil, ErrBadToken
	}
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}
