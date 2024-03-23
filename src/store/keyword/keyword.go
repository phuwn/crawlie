package keyword

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
)

// Store - keyword store interface
type Store interface {
	Get(tx *gorm.DB, name string) (*model.Keyword, error)
	ListByUser(tx *gorm.DB, userID string, limit, offset int, search *string) ([]*model.Keyword, int64, error)
	ListUncrawled(tx *gorm.DB, limit, offset int) ([]*model.Keyword, error)
	BulkInsert(tx *gorm.DB, keywords []*model.Keyword) error
	Update(tx *gorm.DB, keyword *model.Keyword) error
}
