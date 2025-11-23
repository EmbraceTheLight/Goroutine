// demo: 统计 200 篇 RFC 文档中每个字母的出现次数
// 这个 demo 用于测试单线程 / 使用 WaitGroup + Mutex 多线程/ 原子操作+多线程之间的性能对比
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

var mutex sync.Mutex

func countLetters(url string, frequency *[26]int32, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer func() {
		wg.Done()
		resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)

	// 重复搜索 20 遍文档，模拟比较大的文档
	for i := 0; i < 20; i++ {
		for _, b := range body {
			c := strings.ToLower(string(b))
			index := strings.Index(allLetters, c)
			if index >= 0 {
				// 互斥锁
				//mutex.Lock()
				//frequency[index]++
				//mutex.Unlock()

				// 原子操作
				atomic.AddInt32(&frequency[index], 1)
			}
		}
	}
}
func main() {
	var frequency [26]int32
	var wg sync.WaitGroup
	start := time.Now()
	for i := 1000; i <= 1200; i++ {
		wg.Add(1)
		go countLetters(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Done")
	fmt.Printf("Processing took %s\n", elapsed)
	for i, f := range frequency {
		fmt.Printf("%s -> %d\n", string(allLetters[i]), f)
	}
}
