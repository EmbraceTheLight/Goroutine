package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"sync"
)

const (
	screenWidth  = 640
	screenHeight = 360
	boidCount    = 500
	viewRadius   = 13    // 范围。用于寻找半径在 viewRadius 范围内的 boids。
	adjRate      = 0.015 // 调整因子
)

var (
	green = color.RGBA{R: 10, G: 255, B: 50, A: 255}
	boids [boidCount]*Boid

	// 二维 bit map。标记某个位置是否被一个 boid 占据.
	// 两个维度分别表示 x 和 y 坐标，值为 bid。若值为 -1，则表示该位置没有被占据。
	boidMap [screenWidth + 1][screenHeight + 1]int

	lock sync.Mutex
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 打印物体：每个为物体为一个菱形（旋转了 45° 的正方形）
	for _, boid := range boids {
		screen.Set(int(boid.position.x)+1, int(boid.position.y), green)
		screen.Set(int(boid.position.x)-1, int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y)+1, green)
		screen.Set(int(boid.position.x), int(boid.position.y)-1, green)
	}
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	for i, row := range boidMap {
		for j := range row {
			boidMap[i][j] = -1
		}
	}
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boids in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
