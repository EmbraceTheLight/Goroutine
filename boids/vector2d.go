package main

import "math"

type Vector2D struct {
	x float64
	y float64
}

func (v1 *Vector2D) Add(v2 *Vector2D) *Vector2D {
	return &Vector2D{v1.x + v2.x, v1.y + v2.y}
}

func (v1 *Vector2D) Subtract(v2 *Vector2D) *Vector2D {
	return &Vector2D{v1.x - v2.x, v1.y - v2.y}
}

func (v1 *Vector2D) Multiply(v2 *Vector2D) *Vector2D {
	return &Vector2D{v1.x * v2.x, v1.y * v2.y}
}

func (v1 *Vector2D) AddValue(v float64) *Vector2D {
	return &Vector2D{v1.x + v, v1.y + v}
}

func (v1 *Vector2D) SubtractValue(v float64) *Vector2D {
	return &Vector2D{v1.x - v, v1.y - v}
}

func (v1 *Vector2D) MultiplyValue(v float64) *Vector2D {
	return &Vector2D{v1.x * v, v1.y * v}
}

func (v1 *Vector2D) DivisionValue(v float64) *Vector2D {
	return &Vector2D{v1.x / v, v1.y / v}
}

// limit 限制向量的值的大小。
// 如果向量内字段的大小小于lower，则将向量的大小设置为lower。
// 如果向量内字段的大小大于upper，则将向量的大小设置为upper。
func (v1 *Vector2D) limit(lower, upper float64) *Vector2D {
	return &Vector2D{
		x: math.Min(math.Max(v1.x, lower), upper),
		y: math.Min(math.Max(v1.y, lower), upper),
	}
}

func (v1 *Vector2D) Distance(v2 *Vector2D) float64 {
	return math.Sqrt(math.Pow(v1.x-v2.x, 2) + math.Pow(v1.y-v2.y, 2))
}
