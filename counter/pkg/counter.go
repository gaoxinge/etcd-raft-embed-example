package pkg

import "sync"

type Counter struct {
	count    uint64
	mu       sync.RWMutex
	proposeC chan<- struct{}
	commitC  <-chan uint64
}

func NewCounter(proposeC chan<- struct{}, commitC <-chan uint64) *Counter {
	counter := &Counter{
		proposeC: proposeC,
		commitC:  commitC,
	}
	go counter.readCommits()
	return counter
}

func (counter *Counter) Inc() {
	counter.proposeC <- struct{}{}
}

func (counter *Counter) Get() uint64 {
	counter.mu.RLock()
	defer counter.mu.RUnlock()
	return counter.count
}

func (counter *Counter) readCommits() {
	for cnt := range counter.commitC {
		counter.mu.Lock()
		counter.count += cnt
		counter.mu.Unlock()
	}
}
