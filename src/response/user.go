package response

import "github.com/phuwn/crawlie/src/model"

type SignInResponse struct {
	*model.User
	AccessToken string `json:"access_token,omitempty" sql:"-" gorm:"-"`
}
