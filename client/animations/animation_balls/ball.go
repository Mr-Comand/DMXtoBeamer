package animation_balls

import (
	"image/color"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type Balls struct {
	preset_animation.Animation
	balls      []Ball
	peakVolume float32
	dataMap    map[int]int
	friction   float64
}
type Ball struct {
	animation_helpers.Particle
	velocityX float64
	velocityY float64
	friction  *float64
}
type BallsGenerator struct {
}

func NewGeneratorBallsAnimation() *BallsGenerator {
	return &BallsGenerator{}
}
func (g *BallsGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	ballCount := 5
	balls := Balls{
		balls:    make([]Ball, ballCount),
		dataMap:  make(map[int]int),
		friction: 0.0001,
	}
	balls.peakVolume = 0.0
	balls.dataMap = make(map[int]int)
	for i := range balls.balls {
		balls.balls[i].Init(rand.Float64()*50, 0, 500, &balls.friction)
		balls.dataMap[i] = 0
		// Assign random velocities
		balls.balls[i].velocityY = rand.Float64()*8 - 4 // Random Y velocity between -10 and 10
		balls.balls[i].velocityX = rand.Float64()*8 - 4 // Random X velocity between -5 and 5
	}
	return &balls
}

func (b *Ball) Init(size, x, y float64, friction *float64) {
	b.velocityX = 0
	b.velocityY = 0
	b.Particle.Position.X = x
	b.Particle.Position.Y = y
	b.Particle.Color = color.RGBA{255, 0, 0, 255}
	b.Particle.Size = size
	b.friction = friction
}

func (b *Ball) Update() {
	windowWidth := float64(rl.GetScreenWidth())
	windowHeight := float64(rl.GetScreenHeight())

	// Calculate the scaling factors for the window
	scaleX := 0.0
	scaleY := 0.0
	if windowWidth < windowHeight {
		scaleY = (windowHeight - windowWidth) / windowWidth * 1000 / 2
	} else if windowWidth > windowHeight {
		scaleX = (windowWidth - windowHeight) / windowHeight * 1000 / 2
	}
	b.velocityX -= *b.friction * b.velocityX
	b.velocityY -= *b.friction * b.velocityY
	// Update positions
	b.Particle.Position.X += b.velocityX
	b.Particle.Position.Y += b.velocityY
	// Collision detection and response with boundaries
	if b.Particle.Position.X < -scaleX {
		b.Particle.Position.X = -scaleX
		b.velocityX = -b.velocityX // Reverse velocity on collision
	}
	if b.Particle.Position.X > 1000+scaleX {
		b.Particle.Position.X = 1000 + scaleX
		b.velocityX = -b.velocityX // Reverse velocity on collision
	}
	if b.Particle.Position.Y < -scaleY {
		b.Particle.Position.Y = -scaleY
		b.velocityY = -b.velocityY // Reverse velocity on collision
	}
	if b.Particle.Position.Y > 1000+scaleY {
		b.Particle.Position.Y = 1000 + scaleY
		b.velocityY = -b.velocityY // Reverse velocity on collision
	}
}
func (b *Balls) Render(*[]float64) {
	for i := range b.balls {
		b.balls[i].Update()
		b.balls[i].Draw()
		// b.dataMap[i] = int(b.balls[i].Particle.Position.Y)
		if b.dataMap[i] > int(b.peakVolume) {
			b.peakVolume = float32(b.dataMap[i])
		}
	}
}
func (b *Balls) Reset() {

}
func (b *Balls) Draw() {
	for i := range b.balls {
		b.balls[i].Draw()
	}
}
