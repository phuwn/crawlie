package model

import "time"

type KeywordStatus int

const (
	KeywordNeedCrawl KeywordStatus = iota
	KeywordCrawled
)

type Keyword struct {
	ID                 string        `json:"id" gorm:"primaryKey;default:null"`
	Name               string        `json:"name"`
	AdWordsCount       int           `json:"ad_words_count"`
	LinksCount         int           `json:"links_count"`
	SearchResultsCount string        `json:"search_results_count"`
	HtmlCache          string        `json:"html_cache"`
	Status             KeywordStatus `json:"status"`
	LastCrawledAt      *time.Time    `json:"last_crawled_at"`
	UserKeyword        *UserKeyword  `json:"-"`
}
