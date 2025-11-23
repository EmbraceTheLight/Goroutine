package main

import (
	"fmt"
	"time"
)

func waitOnBarrier(name string, timeToSleep int, barrier *Barrier) {
	for {
		for {
			fmt.Println(name, "Running")
			time.Sleep(time.Duration(timeToSleep) * time.Second)
			fmt.Println(name, "is waiting on barrier")
			barrier.Wait()
		}
	}
}
func main() {
	barrier := NewBarrier(2)
	go waitOnBarrier("red", 4, barrier)
	go waitOnBarrier("blue", 10, barrier)
	time.Sleep(100 * time.Second)

}
