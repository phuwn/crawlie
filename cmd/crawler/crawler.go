package main

import (
	"log"

	"github.com/phuwn/crawlie/src/config"
	"github.com/phuwn/crawlie/src/crawler/keyword"
	"github.com/phuwn/crawlie/src/server"
	workerpool "github.com/phuwn/crawlie/src/worker_pool"
	_ "gorm.io/driver/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}

	workerPool, err := workerpool.New(cfg.Crawler.NumberOfWorkers)
	if err != nil {
		log.Fatal(err)
	}

	workloads, err := keyword.LoadUncrawledKeyword(cfg.Crawler.Interval)
	if err != nil {
		log.Fatal(err)
	}

	workerPool.Run(workloads)
}
