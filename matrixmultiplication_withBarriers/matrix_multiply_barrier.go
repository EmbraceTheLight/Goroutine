// 大规模矩阵（n x n 方阵）相乘：使用并发 + 信号量优化
package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	matrixSize = 250
)

var (
	matrixA      = [matrixSize][matrixSize]int{}
	matrixB      = [matrixSize][matrixSize]int{}
	result       = [matrixSize][matrixSize]int{}
	workStart    = NewBarrier(matrixSize + 1) // 多出的一个信号量为主协程准备
	workComplete = NewBarrier(matrixSize + 1)
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	for {
		// 等待主协程生成矩阵 A 和 B
		workStart.Wait()
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
		workComplete.Wait() // 完成一行的计算，等待其他协程完成计算。全部完成后，开启下一轮等待。
	}
}
func main() {
	fmt.Println("Working...")

	// 创建 workOutRow 计算线程。每个 goroutine 扶额计算出新矩阵的一行数据。
	for row := 0; row < matrixSize; row++ {
		go workOutRow(row)
	}

	start := time.Now()

	// 不到 2 s 即可计算完毕。性能与之前使用 waitGroup + 条件变量的效率差不多。但是更直观简洁
	for i := 0; i < 500; i++ {
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		workStart.Wait()    // 生成完毕，此时会广播到所有工作协程中
		workComplete.Wait() // 等待所有工作协程完成计算
	}
	elapsed := time.Since(start)
	fmt.Println("Done...")
	fmt.Printf("Processing took %s\n", elapsed)
}
