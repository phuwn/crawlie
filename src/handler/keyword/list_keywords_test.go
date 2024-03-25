package keyword

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/store"
	keywordMock "github.com/phuwn/crawlie/src/store/keyword/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type listByUserTestCase struct {
	name     string
	limit    string
	offset   string
	q        string
	wantErr  error
	code     int
	response string
}

func TestListByUser(t *testing.T) {
	testcases := []*listByUserTestCase{
		{
			name:     "limit not a number",
			limit:    "test",
			wantErr:  nil,
			code:     400,
			response: `{"error":"limit should be number"}`,
		},
		{
			name:     "offset not a number",
			limit:    "10",
			offset:   "test",
			wantErr:  nil,
			code:     400,
			response: `{"error":"offset should be number"}`,
		},
		{
			name:     "no limit or offset provided",
			wantErr:  nil,
			code:     200,
			response: `{"data":[{"name":"keyword 2","ad_words_count":3,"links_count":4,"search_results_count":"6,000,000","html_cache":"\u003chtml\u003e2\u003c/html\u003e","status":0,"last_crawled_at":null},{"name":"keyword 3","ad_words_count":4,"links_count":5,"search_results_count":"7,000,000","html_cache":"\u003chtml\u003e3\u003c/html\u003e","status":0,"last_crawled_at":null}],"count":2}`,
		},
		{
			name:     "happy case",
			limit:    "10",
			offset:   "2",
			q:        "keyword_test",
			wantErr:  nil,
			code:     200,
			response: `{"data":[{"name":"keyword 1","ad_words_count":1,"links_count":2,"search_results_count":"5,000,000","html_cache":"\u003chtml/\u003e","status":0,"last_crawled_at":null}],"count":1}`,
		},
	}

	keywordStore := keywordMock.NewStore(t)
	keywordStore.Mock.On("ListByUser",
		mock.Anything,
		mock.Anything,
		10,
		2,
		&testcases[3].q,
	).Return([]*model.Keyword{
		{
			Name:               "keyword 1",
			AdWordsCount:       1,
			LinksCount:         2,
			SearchResultsCount: "5,000,000",
			HtmlCache:          "<html/>",
			Status:             model.KeywordNeedCrawl,
		},
	}, int64(1), nil)

	keywordStore.Mock.On("ListByUser",
		mock.Anything,
		mock.Anything,
		50,
		0,
		(*string)(nil),
	).Return([]*model.Keyword{
		{
			Name:               "keyword 2",
			AdWordsCount:       3,
			LinksCount:         4,
			SearchResultsCount: "6,000,000",
			HtmlCache:          "<html>2</html>",
			Status:             model.KeywordNeedCrawl,
		},
		{
			Name:               "keyword 3",
			AdWordsCount:       4,
			LinksCount:         5,
			SearchResultsCount: "7,000,000",
			HtmlCache:          "<html>3</html>",
			Status:             model.KeywordNeedCrawl,
		},
	}, int64(2), nil)
	keywordStore.Mock.On("ListByUser",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil, int64(0), fmt.Errorf("unexpected result occurs"))

	server.SetupTest(nil, nil, &store.Store{
		Keyword: keywordStore,
	}, nil)

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/v1/keywords", nil)
			params := req.URL.Query()
			if tt.limit != "" {
				params.Add("limit", tt.limit)
			}
			if tt.offset != "" {
				params.Add("offset", tt.offset)
			}
			if tt.q != "" {
				params.Add("q", tt.q)
			}
			req.URL.RawQuery = params.Encode()
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("tx", &gorm.DB{})
			err := ListByUser(c)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.code, rec.Code)
			assert.Equal(t, tt.response, rec.Body.String())
		})
	}
}
