package common

import "sync"

// Train 火车结构体
type Train struct {
	Id          int
	TrainLength int
	Front       int // 车头位置
}

// Intersection 交叉口结构体
type Intersection struct {
	Id       int
	Mutex    sync.Mutex
	LockedBy int
}

// Crossing 十字路口结构体
type Crossing struct {
	Position     int // 十字路口位置
	InterStation *Intersection
}
