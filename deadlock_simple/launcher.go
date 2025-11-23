package main

import (
	"github.com/hajimehoshi/ebiten"
	"goroutine/deadlock_simple/arbitrator"
	"goroutine/deadlock_simple/common"
	"goroutine/deadlock_simple/hierarchy"
	"log"
	"sync"
)

var (
	trains        [4]*common.Train
	intersections [4]*common.Intersection
)

const trainLength = 70

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		DrawTracks(screen)
		DrawIntersections(screen)
		DrawTrains(screen)
	}
	return nil
}
func main() {
	for i := 0; i < 4; i++ {
		trains[i] = &common.Train{Id: i, TrainLength: trainLength, Front: 0}
	}

	for i := 0; i < 4; i++ {
		intersections[i] = &common.Intersection{
			Id:       i,
			Mutex:    sync.Mutex{},
			LockedBy: -1,
		}
	}
	hierarchyToBreakDeadLock()
	//arbitratorToBreakDeadLock()
	if err := ebiten.Run(update, 320, 320, 3, "Trains in a box"); err != nil {
		log.Fatal(err)
	}
}

func arbitratorToBreakDeadLock() {
	go arbitrator.MoveTrain(trains[0], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[0]},
		{Position: 175, InterStation: intersections[1]},
	})

	go arbitrator.MoveTrain(trains[1], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[1]},
		{Position: 175, InterStation: intersections[2]},
	})

	go arbitrator.MoveTrain(trains[2], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[2]},
		{Position: 175, InterStation: intersections[3]},
	})

	go arbitrator.MoveTrain(trains[3], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[3]},
		{Position: 175, InterStation: intersections[0]},
	})
}

func hierarchyToBreakDeadLock() {
	go hierarchy.MoveTrain(trains[0], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[0]},
		{Position: 175, InterStation: intersections[1]},
	})

	go hierarchy.MoveTrain(trains[1], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[1]},
		{Position: 175, InterStation: intersections[2]},
	})

	go hierarchy.MoveTrain(trains[2], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[2]},
		{Position: 175, InterStation: intersections[3]},
	})

	go hierarchy.MoveTrain(trains[3], 300, []*common.Crossing{
		{Position: 125, InterStation: intersections[3]},
		{Position: 175, InterStation: intersections[0]},
	})
}
