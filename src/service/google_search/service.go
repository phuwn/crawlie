package googlesearch

import (
	"context"

	"github.com/phuwn/crawlie/src/model"
)

// Service - Google Search Fetching Service
type Service interface {
	Fetch(ctx context.Context, keyword string) ([]byte, error)
	ParseKeywordStatistics(keyword string, htmlPage []byte) (*model.Keyword, error)
}
