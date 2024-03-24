package workerpool

import (
	"log"
	"sync"
)

type WorkLoad interface {
	ID() string
	Work() error
}

func New(numberOfWorkers int) (*WorkerPool, error) {
	return &WorkerPool{
		numberOfWorkers: numberOfWorkers,
	}, nil
}

type WorkerPool struct {
	numberOfWorkers int
}

type counter struct {
	total int
	mu    sync.Mutex
}

func (c *counter) increase() {
	c.mu.Lock()
	c.total++
	c.mu.Unlock()
}

func (c *counter) value() int {
	return c.total
}

func (p *WorkerPool) Run(workloads chan WorkLoad) {
	var (
		wg    sync.WaitGroup
		count counter
		total = len(workloads)
	)

	for i := 0; i < p.numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for workload := range workloads {
				err := workload.Work()
				if err != nil {
					log.Printf("workload %s failed: %s\n", workload.ID(), err.Error())
					continue
				}
				count.increase()
			}
		}()
	}

	wg.Wait()
	succeed := count.value()
	log.Printf("%d records updated!!! %d failed :(\n", succeed, total-succeed)
}
