package executor

import (
	"sync"
)


type ChannelPool struct {
	pool sync.Pool
}


func NewChannelPool() *ChannelPool {
	return &ChannelPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make(chan ResultWithError, 1)
			},
		},
	}
}


func (cp *ChannelPool) Get() chan ResultWithError {
	return cp.pool.Get().(chan ResultWithError)
}


func (cp *ChannelPool) Put(ch chan ResultWithError) {
	select {
	case <-ch:
	default:
	}
	cp.pool.Put(ch)
}


var GlobalChannelPool = NewChannelPool()
