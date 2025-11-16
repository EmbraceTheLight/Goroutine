// 大规模矩阵（n x n 方阵）相乘：使用并发 + 信号量优化
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	matrixSize = 250
)

var (
	matrixA   = [matrixSize][matrixSize]int{}
	matrixB   = [matrixSize][matrixSize]int{}
	result    = [matrixSize][matrixSize]int{}
	rwLock    = sync.RWMutex{}
	cond      = sync.NewCond(rwLock.RLocker())
	waitGroup = sync.WaitGroup{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	rwLock.RLock()
	for {
		waitGroup.Done()
		cond.Wait()
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}
}
func main() {
	fmt.Println("Working...")
	waitGroup.Add(matrixSize)

	// 创建 workOutRow 计算线程。每个 goroutine 扶额计算出新矩阵的一行数据。
	for row := 0; row < matrixSize; row++ {
		go workOutRow(row)
	}

	start := time.Now()
	// 不使用并发，完成 500 次矩阵相乘需要 8.5s；使用并发，是需要不到 2s 即可计算完毕
	for i := 0; i < 500; i++ {
		// 等待 workOutRow 工作线程创建完毕，或者当前矩阵惩乘法计算完毕
		waitGroup.Wait()

		// 上锁，初始化两个矩阵的数据
		rwLock.Lock()
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		waitGroup.Add(matrixSize)
		rwLock.Unlock()
		cond.Broadcast()
	}
	elapsed := time.Since(start)
	fmt.Println("Done...")
	fmt.Printf("Processing took %s\n", elapsed)
}
