package user

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/db"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/util"
	"gorm.io/gorm"
)

func Save(c echo.Context, u *model.User) error {
	cfg := server.GetServerCfg()
	tx, err := db.GetTxFromCtx(c)
	if err != nil {
		return err
	}

	err = cfg.Store().User.Save(tx, u)
	if err != nil {
		return errors.New("save user failed")
	}
	return nil
}

// VerifyGoogleUser - verify user's google auth code, response google info of user
func VerifyGoogleUser(c echo.Context, code, redirectURL string) (*model.User, error) {
	cfg := server.GetServerCfg()
	token, err := cfg.Service().GoogleOauth2.GetToken(c.Request().Context(), code, redirectURL)
	if err != nil {
		return nil, err
	}

	person, err := cfg.Service().GoogleOauth2.GetPerson(c.Request().Context(), token)
	if err != nil {
		return nil, err
	}

	return model.GetUserFromPerson(person)
}

// FirstOrCreate - Get the first record that match user's email or create new user if it doesn't exist
func FirstOrCreate(c echo.Context, u *model.User) error {
	cfg := server.GetServerCfg()
	res, err := cfg.Store().User.GetByEmail(db.Get(), u.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Save(c, u)
		}
		return util.JsonError(c, 404, "failed to get user with email: "+u.Email)
	}

	err = copier.Copy(u, res)
	if err != nil {
		return errors.New("failed to copy value")
	}
	return nil
}

// Get - get user data from the database by the id
func Get(c echo.Context, id string) (*model.User, error) {
	cfg := server.GetServerCfg()
	res, err := cfg.Store().User.Get(db.Get(), id)
	if err != nil {
		return nil, util.JsonError(c, 404, "user not found")
	}
	return res, nil
}
