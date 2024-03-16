package model

import (
	"errors"

	"google.golang.org/api/people/v1"
)

// User data model
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`

	AccessToken *string `json:"access_token,omitempty" sql:"-"`
}

func GetUserFromPerson(person *people.Person) (*User, error) {
	var user User
	if len(person.Names) == 0 {
		return nil, errors.New("missing userinfo.profile scope")
	}
	user.Name = person.Names[0].DisplayName

	if len(person.EmailAddresses) == 0 {
		return nil, errors.New("missing userinfo.email scope")
	}
	user.Email = person.EmailAddresses[0].Value

	if len(person.Photos) == 0 {
		return nil, errors.New("missing userinfo.profile scope")
	}
	user.Avatar = person.Photos[0].Url

	return &user, nil
}
