package keyword

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/phuwn/crawlie/src/model"
	"github.com/stretchr/testify/assert"
)

type extractStatisticTestCase struct {
	name    string
	keyword string
	data    []byte
	result  model.Keyword
	wantErr error
}

func realLifeExampleExtract(lastCrawled *time.Time) (*extractStatisticTestCase, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	found := strings.SplitN(dir, "crawlie", -1)
	if len(found) == 0 {
		return nil, fmt.Errorf("not in crawlie source")
	}

	f, err := os.ReadFile(found[0] + "crawlie/example_page.html")
	if err != nil {
		return nil, err
	}

	// b, err := os.ReadFile(os.P)
	return &extractStatisticTestCase{
		name:    "real life example",
		keyword: "samsung oled tv",
		data:    f,
		result: model.Keyword{
			Name:               "samsung oled tv",
			AdWordsCount:       3,
			LinksCount:         152,
			SearchResultsCount: "167,000,000",
			HtmlCache:          string(f),
			Status:             model.KeywordCrawled,
			LastCrawledAt:      lastCrawled,
		},
	}, nil
}

func TestExtractStatistic(t *testing.T) {
	var testTime = time.Date(2024, time.March, 24, 14, 23, 47, 398324000, time.Local)
	realLifeExample, err := realLifeExampleExtract(&testTime)
	if err != nil {
		t.Error(err)
		return
	}

	testcases := []*extractStatisticTestCase{
		{
			name:    "empty file inserted",
			keyword: "empty",
			result: model.Keyword{
				Name:          "empty",
				Status:        model.KeywordCrawled,
				LastCrawledAt: &testTime,
			},
		},
		{
			name:    "simple data extract",
			keyword: "simple",
			data:    []byte(`<html><div id="result-stats">About 5,000,000 results</div><div data-dtld="Advertiser 1"><div data-dtld="Advertiser 1"/><div data-dtld="Advertiser 3"><a href="link.1">1</a><a href="link.2">2</a><a href="link.3">3</a><a href="link.4">4</a></html>`),
			result: model.Keyword{
				Name:               "simple",
				AdWordsCount:       2,
				LinksCount:         4,
				SearchResultsCount: "5,000,000",
				Status:             model.KeywordCrawled,
				HtmlCache:          `<html><div id="result-stats">About 5,000,000 results</div><div data-dtld="Advertiser 1"><div data-dtld="Advertiser 1"/><div data-dtld="Advertiser 3"><a href="link.1">1</a><a href="link.2">2</a><a href="link.3">3</a><a href="link.4">4</a></html>`,
				LastCrawledAt:      &testTime,
			},
		},
		realLifeExample,
	}

	monkey.Patch(time.Now, func() time.Time {
		return testTime
	})

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			k := KeywordCrawlItem{Keyword: model.Keyword{Name: tt.keyword}}
			resultKeyword, err := k.extractKeywordStatistic(tt.data)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.result, *resultKeyword)
		})
	}
}
