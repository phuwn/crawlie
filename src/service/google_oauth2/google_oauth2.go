package googleoauth2

import (
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type googleOauth2Service struct {
	config *oauth2.Config
}

// NewService - create new google oauth2 service implementation
func NewService(clientID, clientSecret string) Service {
	return &googleOauth2Service{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"profile", "email"},
		},
	}
}

// GetToken - exchange user's code for access_token
func (g *googleOauth2Service) GetToken(ctx context.Context, code, redirectURL string) (*oauth2.Token, error) {
	g.config.RedirectURL = redirectURL
	return g.config.Exchange(ctx, code)
}

// GetPeople - get user's google info
func (g *googleOauth2Service) GetPerson(ctx context.Context, token *oauth2.Token) (*people.Person, error) {
	srv, err := people.NewService(ctx, option.WithHTTPClient(g.config.Client(ctx, token)))
	if err != nil {
		return nil, err
	}

	return srv.People.Get("people/me").
		PersonFields("names,emailAddresses,photos").Do()
}
