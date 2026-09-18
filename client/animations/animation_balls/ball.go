package animation_balls

import (
	"math"
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
	Color           animation_helpers.Color `parameter:"Color,default=#00ff00"`
	BallCount       int                     `parameter:"BallCount,default=5"`
	Shape           int                     `parameter:"Shape,default=0"`
	Hollow          bool                    `parameter:"Hollow,default=false"`
	Size            float64                 `parameter:"Size,default=10"`
	SizeVariation   float64                 `parameter:"SizeVariation,default=50"`
	LineWidth       uint8                   `parameter:"LineWidth,default=0"`
	Rotation        float64                 `parameter:"Rotation,default=0"`
	SoundActive     bool                    `parameter:"SoundActive,default=true"`
	AverageSpeed    float64                 `parameter:"AverageSpeed,default=1000"`
	PushProbability float64                 `parameter:"PushProbability,default=0.005"`
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
		balls.balls[i].Init(rand.Float64()*50, 500, 500, &balls.friction, dynamic.Color)
		balls.dataMap[i] = 0
		// Assign random velocities
		balls.balls[i].velocityY = rand.Float64()*2000 - 1000 // Random Y velocity between -10 and 10
		balls.balls[i].velocityX = rand.Float64()*2000 - 1000 // Random X velocity between -5 and 5
		balls.balls[i].Particle.Shape = animation_helpers.Shape(balls.DynamicConfig.Shape)
		balls.balls[i].Particle.Hollow = balls.DynamicConfig.Hollow
		balls.balls[i].Particle.LineWidth = balls.DynamicConfig.LineWidth
		balls.balls[i].Particle.Rotation = balls.DynamicConfig.Rotation
		balls.balls[i].Particle.Size = balls.DynamicConfig.Size - (rand.Float64()*balls.DynamicConfig.SizeVariation - balls.DynamicConfig.SizeVariation/2)

	}
	return &balls
}

func (b *Ball) Init(size, x, y float64, friction *float64, color animation_helpers.Color) {
	b.velocityX = 0
	b.velocityY = 0
	b.Particle.Position.X = x
	b.Particle.Position.Y = y
	b.Particle.Color = color.RGBA
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

	numBands := len(*data)                   // 2048 frequency bands
	avgSpeed := b.DynamicConfig.AverageSpeed // Maintain this speed

	for i := range b.balls {
		ball := &b.balls[i]

		if b.DynamicConfig.SoundActive && data != nil && numBands > 0 {

			sizeFactor := ball.Particle.Size / 50.0 // Normalize size (50 is an arbitrary max size)

			// Determine frequency range based on ball size
			var startFreq int
			var bandSize int = 10   // Number of frequencies to average
			var sensitivity float64 // Boost factor for higher frequencies

			if sizeFactor > 0.66 {
				startFreq = 5     // Large ball → Low frequencies (Bass)
				sensitivity = 1.0 // Normal bass movement
			} else if sizeFactor > 0.33 {
				startFreq = 200   // Medium ball → Mid frequencies
				sensitivity = 1.5 // Slight boost for mids
			} else {
				startFreq = 1500  // Small ball → High frequencies
				sensitivity = 3.0 // Strong boost for trebles
			}

			// Ensure the frequency range is within bounds
			if startFreq+bandSize >= numBands {
				bandSize = numBands - startFreq - 1
			}

			// Average the values from the selected frequency band
			var frequencySum float64 = 0
			for j := 0; j < bandSize; j++ {
				frequencySum += (*data)[startFreq+j]
			}
			frequencyLevel := (frequencySum / float64(bandSize)) * sensitivity // Apply sensitivity boost

			// Scale movement influence based on size
			movementScale := 2000 * sizeFactor // Bigger balls move more
			velocityBoost := frequencyLevel * movementScale * dt

			// Randomize direction (Up/Down, Left/Right)
			directionY := 1.0
			directionX := 1.0
			if rand.Float64() > 0.5 {
				directionY = -1.0 // Some balls move down instead of up
			}
			if rand.Float64() > 0.5 {
				directionX = -1.0 // Some balls move left instead of right
			}

			// Apply movement with random direction
			ball.velocityY += directionY * velocityBoost
			ball.velocityX += directionX * (velocityBoost * 0.5) // Side movement

		}
		if !b.DynamicConfig.SoundActive {
			// Occasionally apply a random push to keep balls moving
			if rand.Float64() < b.DynamicConfig.PushProbability*dt { // 1% chance per frame
				angle := rand.Float64() * 2 * math.Pi
				pushStrength := 500.0
				ball.velocityX += math.Cos(angle) * pushStrength
				ball.velocityY += math.Sin(angle) * pushStrength
			}
		}

		// Keep speed under control
		speed := math.Sqrt(ball.velocityX*ball.velocityX + ball.velocityY*ball.velocityY)
		if speed > avgSpeed {
			scale := avgSpeed / speed
			ball.velocityX *= scale
			ball.velocityY *= scale
		}
		// Update and draw
		// b.ResolveCollisions() // Call collision handling
		ball.Update(dt)
		ball.Draw()
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
func (b *Balls) HandleCollisions() {
	for i := 0; i < len(b.balls); i++ {
		for j := i + 1; j < len(b.balls); j++ {
			ball1 := &b.balls[i]
			ball2 := &b.balls[j]

			dx := ball2.Particle.Position.X - ball1.Particle.Position.X
			dy := ball2.Particle.Position.Y - ball1.Particle.Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			radiusSum := ball1.Particle.Size/2 + ball2.Particle.Size/2

			if distance < radiusSum && distance > 0 {
				// Normalize the collision vector
				nx := dx / distance
				ny := dy / distance

				// Relative velocity
				vx := ball2.velocityX - ball1.velocityX
				vy := ball2.velocityY - ball1.velocityY

				// Check if they are moving toward each other
				dotProduct := vx*nx + vy*ny
				if dotProduct > 0 {
					continue
				}

				// Apply 1D elastic collision formula along normal
				m1 := ball1.Particle.Size
				m2 := ball2.Particle.Size
				totalMass := m1 + m2

				// Compute new velocities
				ball1.velocityX -= (2 * m2 / totalMass) * dotProduct * nx
				ball1.velocityY -= (2 * m2 / totalMass) * dotProduct * ny
				ball2.velocityX += (2 * m1 / totalMass) * dotProduct * nx
				ball2.velocityY += (2 * m1 / totalMass) * dotProduct * ny

				// Push them apart to prevent overlap
				overlap := radiusSum - distance
				ball1.Particle.Position.X -= nx * (overlap / 2)
				ball1.Particle.Position.Y -= ny * (overlap / 2)
				ball2.Particle.Position.X += nx * (overlap / 2)
				ball2.Particle.Position.Y += ny * (overlap / 2)
			}
		}
	}
}
func (b *Balls) ResolveCollisions() {
	for i := 0; i < len(b.balls); i++ {
		for j := i + 1; j < len(b.balls); j++ {
			ball1, ball2 := &b.balls[i], &b.balls[j]
			dx := ball2.Particle.Position.X - ball1.Particle.Position.X
			dy := ball2.Particle.Position.Y - ball1.Particle.Position.Y
			distance := math.Sqrt(dx*dx + dy*dy)
			radiusSum := ball1.Particle.Size/2 + ball2.Particle.Size/2

			if distance < radiusSum {
				if distance == 0 {
					distance = 0.01
				}

				// Push balls apart to avoid sticking
				overlap := radiusSum - distance
				dxNorm, dyNorm := dx/distance, dy/distance
				ball1.Particle.Position.X -= dxNorm * (overlap / 2)
				ball1.Particle.Position.Y -= dyNorm * (overlap / 2)
				ball2.Particle.Position.X += dxNorm * (overlap / 2)
				ball2.Particle.Position.Y += dyNorm * (overlap / 2)

				// Elastic collision response
				dotProduct := (ball1.velocityX-dxNorm)*dxNorm + (ball1.velocityY-dyNorm)*dyNorm
				ball1.velocityX -= dotProduct * dxNorm
				ball1.velocityY -= dotProduct * dyNorm
				dotProduct = (ball2.velocityX-dxNorm)*dxNorm + (ball2.velocityY-dyNorm)*dyNorm
				ball2.velocityX -= dotProduct * dxNorm
				ball2.velocityY -= dotProduct * dyNorm
			}
		}
	}
}
