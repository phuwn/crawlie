package userkeyword

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userKeywordPGStore struct{}

// NewStore - create new user store implementation with postgreSQL
func NewStore() Store {
	return &userKeywordPGStore{}
}

func (s userKeywordPGStore) BulkInsert(tx *gorm.DB, uks []*model.UserKeyword) error {
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(uks).Error
}
