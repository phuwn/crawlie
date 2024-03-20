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

func NewAuthenticator(cfg config.Authenticator) (*Authenticator, error) {
	var (
		duration = 24 * time.Hour
		err      error
	)
	if cfg.TokenDuration != "" {
		duration, err = time.ParseDuration(cfg.TokenDuration)
		if err != nil {
			return nil, err
		}
	}
	return &Authenticator{jwtSecretKey: cfg.JwtSecretKey, tokenDuration: duration}, nil
}

type Authenticator struct {
	jwtSecretKey  string
	tokenDuration time.Duration
}

// GenerateJWTToken - create an access_token that represents user's session
func (auth *Authenticator) GenerateJWTToken(info *TokenInfo) (string, error) {
	info.ExpiresAt = time.Now().Add(auth.tokenDuration).Unix()
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
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, ErrInvalidSignature
		}
		return nil, ErrBadToken
	}

	if !token.Valid {
		return nil, ErrTokenExpired
	}

	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	return &claims, nil
}
