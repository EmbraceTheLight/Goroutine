package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type SpinLock int32

func (s *SpinLock) Lock() {
	// 自旋锁，使用比较并交换原语
	for !atomic.CompareAndSwapInt32((*int32)(s), 0, 1) {
		// 此时锁被占用。则让出时间片，调度器会切换到其他协程，以便于掌握这个锁的协程进行解锁操作
		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32((*int32)(s), 0)
}

func NewSpinLock() sync.Locker {
	var lock SpinLock
	return &lock
}
func main() {

}
