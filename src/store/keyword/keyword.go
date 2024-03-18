package keyword

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
)

// Store - keyword store interface
type Store interface {
	Get(tx *gorm.DB, name string) (*model.Keyword, error)
	ListByUser(tx *gorm.DB, userID string) ([]*model.Keyword, error)
	BulkInsert(tx *gorm.DB, keywords []*model.Keyword) error
	Save(tx *gorm.DB, user *model.Keyword) error
}
