package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	position *Vector2D
	velocity *Vector2D
	id       int
}

// 计算 boid 的加速度
// 1. 寻找 viewRadius 以内的所有除自身外的 boid，并计算它们的平均速度 avgVelocity
// 2. avgVelocity 减去 boid 的当前速度，得到两者的差值
// 3. 将 差值 乘以调整因子 adjRate，得到 boid 的加速度 accel，
// 加速度的方向 向 viewRadius 以内其他 boid 的平均速度方向靠拢，靠拢幅度取决于调整因子 adjRate
func (b *Boid) calcAcceleration() *Vector2D {
	// upper: 视野范围的右上边界；lower: 视野范围的左下边界
	upper, lower := b.position.AddValue(viewRadius), b.position.AddValue(-viewRadius)

	// 计算在 lower, upper 围成的内切圆区域除自己外的所有 boid 的速度的平均值。
	avgVelocity := &Vector2D{x: 0, y: 0}
	count := 0

	for i := (math.Max(lower.x, 0)); i <= (math.Min(upper.x, screenWidth)); i++ {
		for j := (math.Max(lower.y, 0)); j <= (math.Min(upper.y, screenHeight)); j++ {
			if boidId := boidMap[int(i)][int(j)]; boidId != -1 && boidId != b.id {
				// 判断 boid 是否在视野半径内，若是，则添加速度到平均速度中
				if dist := boids[boidId].position.Distance(b.position); dist < viewRadius {
					avgVelocity = avgVelocity.Add(boids[boidId].velocity)
					count++
				}
			}
		}
	}

	accel := &Vector2D{
		x: 0,
		y: 0,
	}

	if count > 0 {
		fmt.Println("neighbors:", count)
		// 计算平均速度
		avgVelocity = avgVelocity.DivisionValue(float64(count))

		// 2. avgVelocity 减去 boid 的当前速度，得到两者的差值, 若此时 boid 速度加上该加速度，则速度方向与平均速度方向一致
		// 3. 将 差值 乘以调整因子 adjRate，得到 boid 的加速度 accel，此时 boid 速度加上加速度，速度方向趋近平均速度方向，
		accel = avgVelocity.Subtract(b.velocity).MultiplyValue(adjRate)
	}

	return accel
}

func (b *Boid) moveOne() {
	// 更新速度
	b.velocity = b.velocity.Add(b.calcAcceleration()).limit(-1, 1)

	// 记录位置更新前横纵坐标
	oldX, oldY := int(b.position.x), int(b.position.y)

	// 按照既定速度对当前位置进行一次更新, 即：
	// b.position.x = b.position.x + b.velocity.x
	// b.position.y = b.position.y + b.velocity.y
	b.position = b.position.Add(b.velocity)

	// 更新 boidMap
	// 将当前位置标记为 -1，表示已离开该位置
	boidMap[oldX][oldY] = -1

	// 更新新位置的值为 1
	boidMap[int(b.position.x)][int(b.position.y)] = b.id

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
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(bid int) {
	b := Boid{
		position: &Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight}, // 配置随机位置，位于规定的屏幕范围内
		velocity: &Vector2D{rand.Float64()*2 - 1.0, rand.Float64()*2 - 1.0},              // 速度配置为 -1 到 1 之间，
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = bid
	go b.start()
}
