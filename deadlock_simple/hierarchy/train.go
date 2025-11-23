package hierarchy

import (
	"goroutine/deadlock_simple/common"
	"sort"
	"time"
)

func lockIntersectionsInDistance(id, reserveStart, reserveEnd int, crossings []*common.Crossing) {
	var intersectionsToLock []*common.Intersection
	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.InterStation.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.InterStation)
		}
	}

	// 上锁前先对要上锁的交叉口排序，使得每次都是先尝试对编号最小的交叉口上锁，这样做可以打破死锁循环依赖的可能
	sort.Slice(intersectionsToLock, func(i, j int) bool {
		return intersectionsToLock[i].Id < intersectionsToLock[j].Id
	})

	for _, it := range intersectionsToLock {
		it.Mutex.Lock()
		it.LockedBy = id
		time.Sleep(100 * time.Millisecond)
	}
}
func MoveTrain(train *common.Train, distance int, crossings []*common.Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersectionsInDistance(train.Id, crossing.Position, crossing.Position+train.TrainLength, crossings)
			}
			back := train.Front - train.TrainLength
			if back == crossing.Position {
				crossing.InterStation.LockedBy = -1
				crossing.InterStation.Mutex.Unlock()
			}
		}

		time.Sleep(30 * time.Millisecond)
	}
}
