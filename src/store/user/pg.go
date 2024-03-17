package user

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
)

// TODO
type userPGStore struct{}

// NewStore - create new user store implementation why postgreSQL
func NewStore() Store {
	return &userPGStore{}
}

func (s userPGStore) Get(tx *gorm.DB, id string) (*model.User, error) {
	var res model.User
	return &res, tx.Where("id = ?", id).First(&res).Error
}

func (s userPGStore) GetByEmail(tx *gorm.DB, email string) (*model.User, error) {
	u := &model.User{}
	return u, tx.Where("email = ?", email).First(u).Error
}

func (s userPGStore) Save(tx *gorm.DB, user *model.User) error {
	return tx.Save(user).Error
}
