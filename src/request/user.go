package request

// SignInRequest - data form to sign in to auth
type SignInRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}
