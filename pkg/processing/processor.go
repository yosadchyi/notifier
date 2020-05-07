package processing

import (
	"sync"
)

type Processor interface {
	Start(processor func(string) bool)
	Enqueue(item string)
	Stop()
	Wait()
}

type ParallelConfig struct {
	WorkerCount     int
	InputBufferSize int
}

func NewParallelProcessor(config ParallelConfig) Processor {
	return &parallelProcessor{
		ParallelConfig: config,
		c:              make(chan string, config.InputBufferSize),
		wg:             &sync.WaitGroup{},
	}
}

type parallelProcessor struct {
	ParallelConfig
	c  chan string
	wg *sync.WaitGroup
}

func (p *parallelProcessor) Start(processor func(string) bool) {
	p.wg.Add(p.WorkerCount)
	for i := 0; i < p.WorkerCount; i++ {
		go func() {
			for item := range p.c {
				if !processor(item) {
					return
				}
			}
			p.wg.Done()
		}()
	}
}

func (p *parallelProcessor) Enqueue(item string) {
	p.c <- item
}

func (p *parallelProcessor) Stop() {
	close(p.c)
}

func (p *parallelProcessor) Wait() {
	p.wg.Wait()
}
