package service

import (
	"github.com/phuwn/crawlie/src/config"
	googleoauth2 "github.com/phuwn/crawlie/src/service/google_oauth2"
)

// Service - 3rd parties service handling implementation
type Service struct {
	GoogleOauth2 googleoauth2.Service
}

// New - create new service variable
func New(c config.Service) *Service {
	return &Service{
		GoogleOauth2: googleoauth2.NewService(c.GoogleOauth2ClientID, c.GoogleOauth2ClientSecret),
	}
}
