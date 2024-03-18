package userkeyword

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
)

// Store - keyword store interface
type Store interface {
	BulkInsert(tx *gorm.DB, uks []*model.UserKeyword) error
}
