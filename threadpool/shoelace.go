// 使用 shoelace 算法计算多边形面积
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type point2D struct {
	x int
	y int
}

const numberOfThreads = 16

var (
	r         = regexp.MustCompile(`\((\d*),(\d*)\)`)
	waitGroup sync.WaitGroup
)

func findArea(inputChan chan string) {
	defer waitGroup.Done()
	for pointsStr := range inputChan {
		var points []*point2D
		for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, &point2D{x, y})
		}

		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
}
func main() {
	absPath, _ := filepath.Abs("./threadpool/")
	dat, err := os.ReadFile(filepath.Join(absPath, "polygons.txt"))
	if err != nil {
		fmt.Println(err)
	}
	text := string(dat)

	inputChan := make(chan string, 1000)
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inputChan)
	}
	waitGroup.Add(numberOfThreads)
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		//line := "(4,10),(12,8),(10,3),(2,2),(7,5)"
		inputChan <- line
	}
	close(inputChan)

	waitGroup.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s \n", elapsed)
}
