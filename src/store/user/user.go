package user

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
)

// Store - user store interface
type Store interface {
	Get(tx *gorm.DB, id string) (*model.User, error)
	GetByEmail(tx *gorm.DB, email string) (*model.User, error)
	Save(tx *gorm.DB, user *model.User) error
}
