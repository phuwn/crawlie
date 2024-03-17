package handler

import (
	"encoding/json"
	"io"
	"time"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/handler/user"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/util"
)

func userRoutes(g *echo.Group) {
	g.GET("/user/me", getMyInfo)
}

func getMyInfo(c echo.Context) error {
	id := getUserIDFromCtx(c)
	u, err := user.Get(c, id)
	if err != nil {
		return err
	}
	return c.JSON(200, u)
}

// SignInRequest - data form to sign in to auth
type SignInRequest struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

func signIn(c echo.Context) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return util.JsonError(c, 400, "unable to read the request body")
	}

	req := &SignInRequest{}
	err = json.Unmarshal(b, req)
	if err != nil {
		return util.JsonError(c, 400, "wrong sign-in form data")
	}

	u, err := user.VerifyGoogleUser(c, req.Code, req.RedirectURI)
	if err != nil {
		return err
	}

	err = user.FirstOrCreate(c, u)
	if err != nil {
		return err
	}

	jwt, err := generateJWTToken(&TokenInfo{
		User: model.User{
			ID: u.ID,
		},
	}, time.Now().Add(24*time.Hour).Unix())
	if err != nil {
		return err
	}

	u.AccessToken = &jwt
	return c.JSON(200, u)
}
