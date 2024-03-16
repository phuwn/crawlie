package service

import googleoauth2 "github.com/phuwn/crawlie/src/service/google_oauth2"

// Service - 3rd parties service handling implementation
type Service struct {
	Google googleoauth2.Service
}

// New - create new service variable
func New(c *Config) *Service {
	return &Service{
		Google: googleoauth2.NewService(c.GoogleOauth2ClientID, c.GoogleOauth2ClientSecret),
	}
}
