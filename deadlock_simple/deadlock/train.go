package deadlock

import (
	"goroutine/deadlock_simple/common"
	"time"
)

func MoveTrain(train *common.Train, distance int, crossings []*common.Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				crossing.InterStation.Mutex.Lock()
				crossing.InterStation.LockedBy = train.Id
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
