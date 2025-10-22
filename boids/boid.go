package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Boid struct {
	position *Vector2D
	velocity *Vector2D
	id       int
}

func (b *Boid) moveOne() {
	// 按照既定速度对当前位置进行一次更新, 即：
	// b.position.x = b.position.x + b.velocity.x
	// b.position.y = b.position.y + b.velocity.y
	b.position = b.position.Add(b.velocity)

	// 获取下一次 boid 更新后的位置，防止下一次移动时 boid 的位置越过屏幕范围
	next := b.position.Add(b.velocity)

	// 要越过屏幕宽度范围，此时反转速度矢量的 x 轴值
	if next.x < 0 || next.x > screenWidth {
		b.velocity = &Vector2D{
			x: -b.velocity.x,
			y: b.velocity.y,
		}
	}
	// 要越过屏幕高度范围，此时反转速度矢量的 x 轴值
	if next.y < 0 || next.y > screenHeight {
		b.velocity = &Vector2D{
			x: b.velocity.x,
			y: -b.velocity.y,
		}
	}
}

func (b *Boid) start() {
	// 每隔 5 ms 更新一次 boid 的位置坐标
	for {
		b.moveOne()
		if b.id == 1 {
			fmt.Println("moving")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: &Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight}, // 配置随机位置，位于规定的屏幕范围内
		velocity: &Vector2D{rand.Float64()*2 - 1.0, rand.Float64()*2 - 1.0},              // 速度配置为 -1 到 1 之间，
	}
	boids[bid] = &b
	go b.start()
}
