package keyword

import (
	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type keywordPGStore struct{}

// NewStore - create new keyword store implementation with postgreSQL
func NewStore() Store {
	return &keywordPGStore{}
}

func (s keywordPGStore) Get(tx *gorm.DB, name string) (*model.Keyword, error) {
	var res model.Keyword
	return &res, tx.Select("name", "ad_words_count", "links_count", "search_results_count", "status", "html_cache", "last_crawled_at").
		Where("name = ?", name).
		First(&res).Error
}

func (s keywordPGStore) ListByUser(tx *gorm.DB, userID string) ([]*model.Keyword, error) {
	var res []*model.Keyword
	return res, tx.Select("keywords.name", "ad_words_count", "links_count", "search_results_count", "status", "html_cache", "last_crawled_at").
		Joins("left join user_keywords on user_keywords.keyword = keywords.name").
		Joins("left join users on users.id = user_keywords.user_id").
		Where("users.id = ?", userID).
		Find(&res).Error
}

func (s keywordPGStore) BulkInsert(tx *gorm.DB, keywords []*model.Keyword) error {
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"status": model.KeywordNeedCrawl}),
	}).Create(keywords).Error
}

func (s keywordPGStore) Save(tx *gorm.DB, keyword *model.Keyword) error {
	return tx.Save(keyword).Error
}
