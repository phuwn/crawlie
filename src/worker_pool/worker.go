package workerpool

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

type WorkLoad interface {
	ID() string
	Work(client *http.Client, userAgent string) error
}

func New(userAgents []string, interval string) (*WorkerPool, error) {
	var (
		duration = 1 * time.Second
		err      error
	)

	if interval != "" {
		duration, err = time.ParseDuration(interval)
		if err != nil {
			return nil, err
		}
	}

	return &WorkerPool{
		userAgents: userAgents,
		interval:   duration,
	}, nil
}

type WorkerPool struct {
	userAgents []string
	interval   time.Duration
}

type Counter struct {
	mu    sync.Mutex
	total int
}

func (c *Counter) Increase() {
	c.mu.Lock()
	c.total++
	c.mu.Unlock()
}

func (c *Counter) Total() int {
	return c.total
}

func (c *WorkerPool) Run(workloads chan WorkLoad) {
	var wg sync.WaitGroup
	var total = len(workloads)
	var counter Counter
	for _, userAgent := range c.userAgents {
		wg.Add(1)
		ua := userAgent
		go func() {
			jar, _ := cookiejar.New(nil)
			client := &http.Client{Timeout: 5 * time.Second, Jar: jar}
			defer wg.Done()
			for workload := range workloads {
				err := workload.Work(client, ua)
				if err != nil {
					log.Printf("workload %s failed: %s\n", workload.ID(), err.Error())
				}
				counter.Increase()
				time.Sleep(c.interval)
			}
		}()
	}

	wg.Wait()
	succeed := counter.Total()
	log.Printf("%d records updated!!! %d failed :(\n", succeed, total-succeed)
}
