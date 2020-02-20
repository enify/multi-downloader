package worker

import (
	"sync"
)

// Pool parallel support
type Pool struct {
	wg sync.WaitGroup

	MaxGoroutines int

	workchan chan Work
	Errchan  chan error
}

// NewPool create an new working pool
func NewPool(maxGoroutines, bufferSize int) *Pool {
	p := Pool{
		MaxGoroutines: maxGoroutines,
		workchan:      make(chan Work, bufferSize),
		Errchan:       make(chan error, bufferSize),
	}
	p.wg.Add(maxGoroutines)

	return &p
}

// Run start working pool
func (p *Pool) Run() {
	for i := 0; i < p.MaxGoroutines; i++ {
		go func() {
			for work := range p.workchan {
				err := work.Do()
				if err != nil {
					select {
					case p.Errchan <- err:
					default:
					}
				}
			}
			p.wg.Done()
		}()
	}
}

// ShutDown waits for all goroutines to be shutdown
func (p *Pool) ShutDown() {
	close(p.workchan)
	close(p.Errchan)
	p.wg.Wait()
}

// Submit add new work to working pool
func (p *Pool) Submit(w Work) {
	p.workchan <- w
}
