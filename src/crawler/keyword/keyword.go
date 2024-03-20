package keyword

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/server"
	"github.com/phuwn/crawlie/src/util"
	workerpool "github.com/phuwn/crawlie/src/worker_pool"
	"golang.org/x/net/html"
)

var (
	searchResultsReg = regexp.MustCompile("About ([0-9,]+) results")
)

func LoadUncrawledKeyword() (chan workerpool.WorkLoad, error) {
	srv := server.Get()
	keywords, err := srv.Store().Keyword.ListUncrawled(srv.DB().DB(), 50, 0)
	if err != nil {
		return nil, err
	}

	workloads := make(chan workerpool.WorkLoad, len(keywords))
	for _, keyword := range keywords {
		workloads <- &KeywordCrawlItem{Keyword: *keyword}
	}
	close(workloads)
	return workloads, nil
}

type KeywordCrawlItem struct {
	model.Keyword
}

func (k KeywordCrawlItem) ID() string {
	return k.Name
}

func (k KeywordCrawlItem) Work(client *http.Client, userAgent string) error {
	srv := server.Get()
	b, err := k.fetchSearchPage(client, userAgent)
	if err != nil {
		return err
	}

	keyword, err := k.extractKeywordStatistic(b)
	if err != nil {
		return err
	}

	return srv.Store().Keyword.Save(srv.DB().DB(), keyword)
}

func (k KeywordCrawlItem) fetchSearchPage(client *http.Client, userAgent string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://www.google.com/search", nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{
		"q":  {k.Name},
		"oq": {k.Name},
		"hl": {"en"},
		"gl": {"en,"},
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (k KeywordCrawlItem) extractKeywordStatistic(b []byte) (*model.Keyword, error) {
	doc, err := html.Parse(bytes.NewReader(b))
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
		Name:               k.Name,
		AdWordsCount:       len(advertiserMap),
		LinksCount:         linksCount,
		SearchResultsCount: searchResultsCount,
		HtmlCache:          string(b),
		Status:             model.KeywordCrawled,
		LastCrawledAt:      &now,
	}, nil
}
