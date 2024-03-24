package keyword

import (
	"context"
	"log"
	"time"

	"github.com/phuwn/crawlie/src/model"
	"github.com/phuwn/crawlie/src/server"
	workerpool "github.com/phuwn/crawlie/src/worker_pool"
)

func LoadUncrawledKeyword(intervalStr string) (chan workerpool.WorkLoad, error) {
	var (
		interval = 2 * time.Second
		err      error
	)

	if intervalStr != "" {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			return nil, err
		}
	}

	srv := server.Get()
	keywords, err := srv.Store().Keyword.ListUncrawled(srv.DB().DB(), 50, 0)
	if err != nil {
		return nil, err
	}
	log.Printf("%d uncrawled keywords found.\n", len(keywords))

	workloads := make(chan workerpool.WorkLoad, len(keywords))
	for _, keyword := range keywords {
		workloads <- &KeywordCrawlItem{Keyword: *keyword, interval: interval}
	}
	close(workloads)
	return workloads, nil
}

type KeywordCrawlItem struct {
	model.Keyword
	interval time.Duration
}

func (k KeywordCrawlItem) ID() string {
	return k.Name
}

func (k KeywordCrawlItem) Work() error {
	time.Sleep(k.interval)
	srv := server.Get()
	log.Println("Fetching keyword ", k.Name)
	b, err := srv.Service().GoogleSearch.Fetch(context.Background(), k.Name)
	if err != nil {
		return err
	}

	keyword, err := srv.Service().GoogleSearch.ParseKeywordStatistics(k.Name, b)
	if err != nil {
		return err
	}

	return srv.Store().Keyword.Save(srv.DB().DB(), keyword)
}
