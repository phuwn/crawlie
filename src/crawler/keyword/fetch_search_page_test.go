package keyword

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/phuwn/crawlie/src/model"
	"github.com/stretchr/testify/assert"
)

type fetchSearchPageTestCase struct {
	name         string
	keyword      string
	userAgent    string
	responseData string
	wantErr      error
}

func TestFetchSearchPage(t *testing.T) {
	testcases := []*fetchSearchPageTestCase{
		{
			name:         "happy case",
			keyword:      "200",
			userAgent:    "Mozilla/5.0",
			responseData: "<html>succeed</html>",
			wantErr:      nil,
		},
		{
			name:    "unexpected error occurs",
			keyword: "429",
			wantErr: fmt.Errorf("request failed with status code 429"),
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://www.google.com/search",
		func(req *http.Request) (*http.Response, error) {
			params := req.URL.Query()
			headers := req.Header
			if params.Get("q") == "200" && headers.Get("User-Agent") == "Mozilla/5.0" {
				return httpmock.NewStringResponse(200, "<html>succeed</html>"), nil
			}
			return httpmock.NewStringResponse(429, ""), nil
		})

	client := &http.Client{Timeout: 5 * time.Second}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			k := KeywordCrawlItem{Keyword: model.Keyword{Name: tt.keyword}}
			b, err := k.fetchSearchPage(client, tt.userAgent)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.responseData, string(b))
		})
	}
}
