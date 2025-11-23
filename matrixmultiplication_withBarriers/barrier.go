package main

import "sync"

type Barrier struct {
	total int // 屏障阻塞 goroutine 的数量
	count int // 屏障还需要多少个 goroutine 才能解除阻塞
	mutex *sync.Mutex
	cond  *sync.Cond
}

func NewBarrier(size int) *Barrier {
	lockToUse := &sync.Mutex{}
	condToUse := sync.NewCond(lockToUse)
	return &Barrier{
		total: size,
		count: size,
		mutex: lockToUse,
		cond:  condToUse,
	}
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	b.count -= 1
	if b.count == 0 {
		b.count = b.total
		b.cond.Broadcast()
	} else {
		b.cond.Wait()
	}
	b.mutex.Unlock()
}
