// 转账 demo: 测试自旋锁与互斥锁的性能差异
package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	totalAccounts  = 500
	maxAmountMoved = 10
	initialMoney   = 1000
	threads        = 4
)

// 模拟转账事务
func performMovements(ledger *[totalAccounts]int64, locks *[totalAccounts]sync.Locker, totalTrans *int64) {
	for {
		// 1. 随机选择两个账户进行转账，从 accountA 账户转账 amountToMove 到 accountB 账户。
		accountA := rand.Intn(totalAccounts)
		accountB := rand.Intn(totalAccounts)
		// 两个账户不能相同
		for accountA == accountB {
			accountB = rand.Intn(totalAccounts)
		}

		// 2. 随机选择转账金额，不超过 maxAmountMoved
		amountToMove := rand.Int63n(maxAmountMoved)
		toLock := []int{accountA, accountB}

		// 3 排序账户锁，保证锁定的账户顺序一致，避免死锁。见 ./ledger_example_deadlock.png
		sort.Ints(toLock)
		locks[toLock[0]].Lock()
		locks[toLock[1]].Lock()
		atomic.AddInt64(&ledger[accountA], -amountToMove)
		atomic.AddInt64(&ledger[accountB], amountToMove)
		atomic.AddInt64(totalTrans, 1) // 转账事务 + 1

		locks[toLock[0]].Unlock()
		locks[toLock[1]].Unlock()

	}
}
func main() {
	fmt.Println("Total Accounts:", totalAccounts, "Total Accounts:", totalAccounts, "Threads:", threads)
	var ledger [totalAccounts]int64
	var locks [totalAccounts]sync.Locker
	var totalTrans int64
	for i := 0; i < totalAccounts; i++ {
		ledger[i] = initialMoney
		locks[i] = NewSpinLock()
		//locks[i] = &sync.Mutex{}
	}

	// 启动 4 个协程，模拟并发转账场景
	for i := 0; i < threads; i++ {
		go performMovements(&ledger, &locks, &totalTrans)
	}

	// 统计总金额 ---- 它应当保持不变
	for {
		time.Sleep(2 * time.Second)
		var sum int64
		for i := 0; i < totalAccounts; i++ {
			locks[i].Lock()
		}
		for i := 0; i < totalAccounts; i++ {
			sum += ledger[i]
		}
		for i := 0; i < totalAccounts; i++ {
			locks[i].Unlock()
		}
		fmt.Println(totalTrans, sum)
	}
}
