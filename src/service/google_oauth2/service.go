package googleoauth2

import (
	"context"

	"golang.org/x/oauth2"
	"google.golang.org/api/people/v1"
)

// Service - Google Oauth 2.0 service
type Service interface {
	GetToken(ctx context.Context, code, redirectURL string) (*oauth2.Token, error)
	GetPerson(ctx context.Context, token *oauth2.Token) (*people.Person, error)
}
