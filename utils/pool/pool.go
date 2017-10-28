package pool

import (
	"errors"
	"sync"
	//"time"
)

type Pool interface {
	Put(interface{})
	Get() (interface{}, error)
	Size() int64
	Usable() int64
}

func NewPool(size int) (*Pool, error) {
	if size < 1 {
		return nil, errors.New("the pool size must be greater than 1")
	}
	p := &pool{size: size}
	p.used = 0
	p.job = make(chan interface{}, size)

	return p
}

type pool struct {
	size int64 // 池的总容量。
	used int64 // 池的使用量
	job  chan interface{}
	mu   sync.Mutex
}

func (p *pool) Get() (interface{}, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	job, ok := <-p.job
	if !ok {
		return nil, errors.New("get job error")
	}

	p.used--

	return job, nil

}

func (p *pool) Put(job interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.job <- job
	p.used++
}

func (p *pool) Size() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.size
}

func (p *pool) Usable() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()

	usable := p.size - p.used

	return usable
}
