package googlesearch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"time"

	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/util"
	"golang.org/x/net/html"
)

type googleSearchService struct {
	userAgents []string
}

// NewService - create new google search service implementation
func NewService(userAgents []string) Service {
	return &googleSearchService{userAgents: userAgents}
}

// Fetch - retrieve the first result page of the search keyword
func (g *googleSearchService) Fetch(ctx context.Context, keyword string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com/search", nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{
		"q":  {keyword},
		"oq": {keyword},
		"hl": {"en"},
		"gl": {"en,"},
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("User-Agent", g.userAgents[util.RandRange(0, len(g.userAgents))])

	jar, _ := cookiejar.New(nil)
	client := http.Client{Timeout: 5 * time.Second, Jar: jar}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

var (
	searchResultsReg = regexp.MustCompile("About ([0-9,]+) results")
)

func (g *googleSearchService) ParseKeywordStatistics(keyword string, htmlPage []byte) (*model.Keyword, error) {
	doc, err := html.Parse(bytes.NewReader(htmlPage))
	if err != nil {
		return nil, err
	}

	var (
		linksCount         int
		searchResultsCount string
		traverse           func(n *html.Node)
	)

	advertiserMap := make(map[string]bool)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			id, found := util.FindHTMLAttribute(n, "id")
			if found && id == "result-stats" {
				fc := n.FirstChild
				if fc != nil {
					matches := searchResultsReg.FindStringSubmatch(fc.Data)
					if len(matches) > 1 {
						searchResultsCount = matches[1]
					}
				}
			}

			advertiser, found := util.FindHTMLAttribute(n, "data-dtld")
			if found {
				advertiserMap[advertiser] = true
			}

			if n.Data == "a" {
				_, found := util.FindHTMLAttribute(n, "href")
				if found {
					linksCount++
				}
			}

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	now := time.Now()
	return &model.Keyword{
		Name:               keyword,
		AdWordsCount:       len(advertiserMap),
		LinksCount:         linksCount,
		SearchResultsCount: searchResultsCount,
		HtmlCache:          string(htmlPage),
		Status:             model.KeywordCrawled,
		LastCrawledAt:      &now,
	}, nil
}
