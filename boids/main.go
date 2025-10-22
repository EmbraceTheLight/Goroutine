package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 360
	boidCount    = 500
)

var (
	green = color.RGBA{R: 10, G: 255, B: 50, A: 255}
	boids [boidCount]*Boid
)

type Game struct{}

//func (g *Game) Update(screen *ebiten.Image) error {
//	return nil
//}

func (g *Game) Update(screen *ebiten.Image) error {
	// 打印物体：每个为物体为一个菱形（旋转了 45° 的正方形）
	for _, boid := range boids {
		screen.Set(int(boid.position.x)+1, int(boid.position.y), green)
		screen.Set(int(boid.position.x)-1, int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y)+1, green)
		screen.Set(int(boid.position.x), int(boid.position.y)-1, green)
	}
	return nil
}

func (g *Game) Layout(_, _ int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Boids in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
