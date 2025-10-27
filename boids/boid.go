package main

import (
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

	// 计算在 lower, upper 围成的内切圆区域除自己外的所有 boid 的位置的平均值、速度的平均值和分离加速度的平均值。
	avgPosition, avgVelocity, separation := &Vector2D{x: 0, y: 0}, &Vector2D{x: 0, y: 0}, &Vector2D{x: 0, y: 0}
	count := 0.0

	rwLock.RLock()
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if boidId := boidMap[int(i)][int(j)]; boidId != -1 && boidId != b.id {
				// 判断 boid 是否在视野半径内，若是，则添加速度到平均速度中
				if dist := boids[boidId].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[boidId].velocity)
					avgPosition = avgPosition.Add(boids[boidId].position)
					separation = separation.Add(b.position.Subtract(boids[boidId].position)).DivisionValue(dist)
				}
			}
		}
	}
	rwLock.RUnlock()

	accel := &Vector2D{
		x: b.borderBounce(b.position.x, screenWidth),
		y: b.borderBounce(b.position.y, screenHeight),
	}

	if count > 0 {
		// 计算平均速度、平均位置
		avgVelocity = avgVelocity.DivisionValue(count)
		avgPosition = avgPosition.DivisionValue(count)

		// 2. avgVelocity 减去 boid 的当前速度，得到两者的差值, 若此时 boid 速度加上该加速度，则速度方向与平均速度方向一致
		// 3. 将 差值 乘以调整因子 adjRate，得到 boid 的加速度 accel，此时 boid 速度加上加速度，速度方向趋近平均速度方向，
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyValue(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyValue(adjRate)
		accelSeparation := separation.MultiplyValue(adjRate)
		accel = accel.
			Add(accelAlignment).
			Add(accelCohesion).
			Add(accelSeparation)
	}

	return accel
}

// borderBounce 边界反弹
func (b *Boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius { // 在左下角
		return 1 / pos

		// 在右上角
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}
	return 0
}

func (b *Boid) moveOne() {
	acceleration := b.calcAcceleration()
	rwLock.Lock()

	// 更新速度
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1)

	/// 更新 boidMap
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	// 更新新位置的值为 1
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	rwLock.Unlock()
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
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = bid
	go b.start()
}
