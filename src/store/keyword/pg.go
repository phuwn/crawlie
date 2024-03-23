package keyword

import (
	"time"

	"github.com/phuwn/crawlie/src/model"
	"gorm.io/gorm"
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

func (s keywordPGStore) ListByUser(tx *gorm.DB, userID string, limit, offset int, search *string) ([]*model.Keyword, int64, error) {
	var res []*model.Keyword
	var count int64

	tx1 := tx.Table("keywords").
		Joins("left join user_keywords on user_keywords.keyword_id = keywords.id").
		Where("user_keywords.user_id = ?", userID)

	if search != nil {
		tx1 = tx1.Where("keywords.name LIKE ?", "%"+*search+"%")
	}
	err := tx1.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	tx2 := tx.Select("keywords.name", "ad_words_count", "links_count", "search_results_count", "status", "html_cache", "last_crawled_at").
		Joins("left join user_keywords on user_keywords.keyword_id = keywords.id").
		Where("user_keywords.user_id = ?", userID)

	if search != nil {
		tx2 = tx2.Where("keywords.name LIKE ?", "%"+*search+"%")
	}

	return res, count, tx2.Limit(limit).Offset(offset).Find(&res).Error
}

func (s keywordPGStore) ListUncrawled(tx *gorm.DB, limit, offset int) ([]*model.Keyword, error) {
	var res []*model.Keyword
	return res, tx.Select("name", "ad_words_count", "links_count", "search_results_count", "status", "html_cache", "last_crawled_at").
		Where("status = ?", model.KeywordNeedCrawl).Where("last_crawled_at is null or last_crawled_at < ?", time.Now().Add(-10*time.Minute)).
		Limit(limit).Offset(offset).Find(&res).Error
}

func (s keywordPGStore) BulkInsert(tx *gorm.DB, keywords []*model.Keyword) error {
	return tx.Create(keywords).Error
}

func (s keywordPGStore) Update(tx *gorm.DB, keyword *model.Keyword) error {
	return tx.Omit("UserKeyword").Save(keyword).Error
}
