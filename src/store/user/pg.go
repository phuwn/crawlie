package user

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userPGStore struct{}

// NewStore - create new user store implementation with postgreSQL
func NewStore() Store {
	return &userPGStore{}
}

func (s userPGStore) Get(tx *gorm.DB, id string) (*model.User, error) {
	var res model.User
	return &res, tx.Select("id", "name", "email", "avatar").Where("id = ?", id).First(&res).Error
}

func (s userPGStore) GetByEmail(tx *gorm.DB, email string) (*model.User, error) {
	u := &model.User{}
	return u, tx.Select("id", "name", "email", "avatar").Where("email = ?", email).First(u).Error
}

func (s userPGStore) Save(tx *gorm.DB, user *model.User) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "avatar"}),
	}).Create(user).Error
}
