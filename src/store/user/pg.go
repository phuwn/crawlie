package user

import (
	"context"

	"github.com/phuwn/crawlie/src/model"
)

// TODO
type userPGStore struct{}

// NewStore - create new user store implementation why postgreSQL
func NewStore() Store {
	return &userPGStore{}
}

func (s userPGStore) Get(c context.Context, id string) (*model.User, error) {
	// tx := db.GetTxFromCtx(c)
	// var res model.User
	// return &res, tx.Where("id = ?", id).First(&res).Error
	return nil, nil
}

func (s userPGStore) GetByEmail(c context.Context, email string) (*model.User, error) {
	// tx := db.GetTxFromCtx(c)
	// u := &model.User{}
	// return u, tx.Where("email = ?", email).First(u).Error
	return nil, nil
}

func (s userPGStore) Save(c context.Context, user *model.User) error {
	// tx := db.GetTxFromCtx(c)
	// return tx.Save(user).Error
	return nil
}
