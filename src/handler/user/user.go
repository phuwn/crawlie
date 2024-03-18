package user

import (
	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/auth"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/request"
	"github.com/phuwn/crawlie/src/response"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/util"
)

// verifyGoogleUser - verify user's google auth code, response google info of user
func verifyGoogleUser(c echo.Context, code, redirectURL string) (*model.User, error) {
	srv := server.Get()
	token, err := srv.Service().GoogleOauth2.GetToken(c.Request().Context(), code, redirectURL)
	if err != nil {
		return nil, err
	}

	person, err := srv.Service().GoogleOauth2.GetPerson(c.Request().Context(), token)
	if err != nil {
		return nil, err
	}

	return model.GetUserFromPerson(person)
}

// Get - get user data by ID
func Get(c echo.Context) error {
	var (
		id = util.GetUserIDFromCtx(c)
		tx = util.GetTxFromCtx(c)
	)
	u, err := server.Get().Store().User.Get(tx, id)
	if err != nil {
		return util.JsonError(c, 404, "user not found")
	}
	return c.JSON(200, u)
}

func SignIn(c echo.Context) error {
	req := &request.SignInRequest{}
	err := util.JsonParse(c, req)
	if err != nil {
		return err
	}

	u, err := verifyGoogleUser(c, req.Code, req.RedirectURI)
	if err != nil {
		return err
	}

	var (
		tx  = util.GetTxFromCtx(c)
		srv = server.Get()
	)

	err = srv.Store().User.Save(tx, u)
	if err != nil {
		return err
	}

	jwt, err := srv.Auth().GenerateJWTToken(&auth.TokenInfo{
		User: *u,
	})
	if err != nil {
		return err
	}

	return c.JSON(200, &response.SignInResponse{
		User:        u,
		AccessToken: jwt,
	})
}
