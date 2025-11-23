// 统计 200 篇 RFC 文档中每个字母的出现次数
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency *[26]int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		index := strings.Index(allLetters, c)
		if index >= 0 {
			frequency[index]++
		}
	}
}
func main() {
	var frequency [26]int
	start := time.Now()
	for i := 1000; i <= 1200; i++ {
		fmt.Println("Fetching RFC", fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i))
		countLetters(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency)
	}
	elapsed := time.Since(start)
	fmt.Println("Done")
	fmt.Printf("Processing took %s\n", elapsed)
	for i, f := range frequency {
		fmt.Printf("%s -> %d\n", string(allLetters[i]), f)
	}
}
