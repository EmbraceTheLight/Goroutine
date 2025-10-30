package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	matches   []string
	waitgroup = sync.WaitGroup{}
	limit     = make(chan struct{}, 20)
	lock      = sync.Mutex{}
)

func fileSearch(root string, filename string) {
	//fmt.Println("Searching in", root)
	defer func() {
		waitgroup.Done()
		<-limit
	}()

	limit <- struct{}{}
	files, _ := os.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			lock.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}
		if file.IsDir() {
			waitgroup.Add(1)
			go fileSearch(filepath.Join(root, file.Name()), filename)
		}
	}
}
func main() {
	t := time.Now()
	waitgroup.Add(1)
	go fileSearch("C:/", "Readme.md")
	waitgroup.Wait()
	for _, file := range matches {
		fmt.Println("Matched", file)
	}
	fmt.Println("Time taken:", time.Since(t))
}
