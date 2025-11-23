package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	money int32 = 100
)

func stingy() {
	for i := 1; i <= 100000; i++ {
		atomic.AddInt32(&money, 10)
	}
	fmt.Println("Stingy Done")
}

func spendy() {
	for i := 1; i <= 100000; i++ {
		atomic.AddInt32(&money, -10)
	}
	fmt.Println("Spendy Done")
}

func main() {
	go stingy()
	go spendy()

	time.Sleep(3000 * time.Millisecond)
	fmt.Println("Money:", money)
}
