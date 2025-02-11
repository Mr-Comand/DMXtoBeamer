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
	balls         []Ball
	peakVolume    float32
	dataMap       map[int]int
	friction      float64
	DynamicConfig *DynamicConfig
}
type DynamicConfig struct {
	Color     animation_helpers.Color `parameter:"Color,default=ff0000"` //TODO
	BallCount int                     `parameter:"BallCount,default=5"`
	Shape     int                     `parameter:"Shape,default=0"`
	Hollow    bool                    `parameter:"Hollow,default=false"`
	LineWidth uint8                   `parameter:"LineWidth,default=0"`
	Rotation  float64                 `parameter:"Rotation,default=0"`
	// TODO: Sound Active
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
func (g *BallsGenerator) Unload() {
}
func (g *BallsGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	dynamic := &DynamicConfig{}
	preset_animation.Parse(config, dynamic)
	ballCount := dynamic.BallCount
	balls := Balls{
		balls:         make([]Ball, ballCount),
		dataMap:       make(map[int]int),
		friction:      0.1,
		DynamicConfig: dynamic,
	}
	balls.peakVolume = 0.0
	balls.dataMap = make(map[int]int)
	for i := range balls.balls {
		balls.balls[i].Init(rand.Float64()*50, 500, 500, &balls.friction)
		balls.dataMap[i] = 0
		// Assign random velocities
		balls.balls[i].velocityY = rand.Float64()*2000 - 1000 // Random Y velocity between -10 and 10
		balls.balls[i].velocityX = rand.Float64()*2000 - 1000 // Random X velocity between -5 and 5
		balls.balls[i].Particle.Shape = animation_helpers.Shape(balls.DynamicConfig.Shape)
		balls.balls[i].Particle.Hollow = balls.DynamicConfig.Hollow
		balls.balls[i].Particle.LineWidth = balls.DynamicConfig.LineWidth
		balls.balls[i].Particle.Rotation = balls.DynamicConfig.Rotation
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

func (b *Ball) Update(dt float64) {
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
	b.velocityX -= *b.friction * dt * b.velocityX
	b.velocityY -= *b.friction * dt * b.velocityY
	// Update positions

	b.Particle.Position.X += b.velocityX * (dt)
	b.Particle.Position.Y += b.velocityY * (dt)
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
func (b *Balls) Render(data *[]float64, dt float64) {
	for i := range b.balls {
		b.balls[i].Update(dt)
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
func (a *Balls) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	for i := range a.balls {
		a.balls[i].Particle.Shape = animation_helpers.Shape(a.DynamicConfig.Shape)
		a.balls[i].Particle.Hollow = a.DynamicConfig.Hollow
		a.balls[i].Particle.LineWidth = a.DynamicConfig.LineWidth
		a.balls[i].Particle.Rotation = a.DynamicConfig.Rotation
	}
}
