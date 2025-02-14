package animation_falling_shape

import (
	"image/color"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type FallingShape struct {
	preset_animation.Animation
	FallObjects      []FallObject
	DynamicConfig    *DynamicConfig
	WindStrength     float64 // Current wind strength
	WindDirection    int     // Wind direction: -1 (left) or 1 (right)
	WindGustCooldown float64 // Cooldown timer for wind gusts
}

type DynamicConfig struct {
	Color            animation_helpers.Color `parameter:"Color,default=ff0000"`
	Shape            animation_helpers.Shape `parameter:"Shape,default=1"`
	Hollow           bool                    `parameter:"Hollow,default=false"`
	Size             float64                 `parameter:"Size,default=50"`
	SizeVariation    float64                 `parameter:"SizeVariation,default=10"`
	LineWidth        uint8                   `parameter:"LineWidth,default=50"`
	Rotation         float64                 `parameter:"Rotation,default=0"`
	SpinSpeed        float64                 `parameter:"SpinSpeed,default=0.5"`
	Swing            bool                    `parameter:"Swing,default=true"`
	Gravity          float64                 `parameter:"Gravity,default=100"`
	ParticleCount    int                     `parameter:"ParticleCount,default=10"`
	WindgutsStrength float32                 `parameter:"WindgutsStrength,default=0"`
}
type FallObject struct {
	animation_helpers.Particle
	Z     float32
	Speed float64
	Phase float64
}
type FallingShapeGenerator struct{}

func NewFallingShapeGenerator() *FallingShapeGenerator {
	return &FallingShapeGenerator{}
}

func (g *FallingShapeGenerator) Unload() {}

func (g *FallingShapeGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	shape := FallingShape{
		DynamicConfig: &DynamicConfig{},
	}
	shape.FallObjects = make([]FallObject, shape.DynamicConfig.ParticleCount)
	shape.Configure(config)
	return &shape
}

func (a *FallingShape) Render(data *[]float64, dt float64) {
	a.Update(dt)
	for i := range a.FallObjects {
		a.FallObjects[i].Draw()
	}
}

func (a *FallingShape) Update(dt float64) {
	windowWidth := float64(rl.GetScreenWidth())
	windowHeight := float64(rl.GetScreenHeight())
	scaleX := 0.0
	scaleY := 0.0
	if windowWidth < windowHeight {
		scaleY = (windowHeight - windowWidth) / windowWidth * 1000 / 2
	} else if windowWidth > windowHeight {
		scaleX = (windowWidth - windowHeight) / windowHeight * 1000 / 2
	}
	time := rl.GetTime()

	// Update wind gust cooldown
	if a.WindGustCooldown > 0 {
		a.WindGustCooldown -= dt
	}

	// Trigger wind gusts periodically
	if a.WindGustCooldown <= 0 && rand.Float64() < 0.01 { // 1% chance to trigger a wind gust
		a.TriggerWindGust()
		a.WindGustCooldown = 5 // Set cooldown to 5 seconds
	}

	// Gradually reduce wind strength over time
	if a.WindStrength > 0 {
		a.WindStrength -= 0.1 * dt
		if a.WindStrength < 0 {
			a.WindStrength = 0
		}
	}

	// Update particle positions
	for i, p := range a.FallObjects {
		// Apply gravity
		a.FallObjects[i].Position.Y += p.Speed * dt

		// Apply wind effect
		a.FallObjects[i].Position.X += a.WindStrength * float64(a.WindDirection) * dt

		// Apply swing effect
		if a.DynamicConfig.Swing {
			a.FallObjects[i].Position.X += p.randomSwing(time) * dt * 100
		}

		// Rotate if FreeSpin is enabled
		if a.DynamicConfig.SpinSpeed >= 0 {
			a.FallObjects[i].Rotation += math.Cos(time+p.Phase) * dt * 100 * a.DynamicConfig.SpinSpeed
		}

		// Reset position if it falls out of bounds
		if a.FallObjects[i].Position.Y > 1000+scaleY+a.FallObjects[i].Size || a.FallObjects[i].Position.X > 1500+scaleX+a.FallObjects[i].Size || a.FallObjects[i].Position.X < -500-scaleX-a.FallObjects[i].Size {
			a.FallObjects[i].Reset(a.DynamicConfig, scaleX, scaleY)
		}
	}
}

func (a *FallingShape) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	if a.DynamicConfig.ParticleCount > len(a.FallObjects) {
		windowWidth := float64(rl.GetScreenWidth())
		windowHeight := float64(rl.GetScreenHeight())
		scaleX := 0.0
		scaleY := 0.0
		if windowWidth < windowHeight {
			scaleY = (windowHeight - windowWidth) / windowWidth * 1000 / 2
		} else if windowWidth > windowHeight {
			scaleX = (windowWidth - windowHeight) / windowHeight * 1000 / 2
		}
		for i := 0; i < a.DynamicConfig.ParticleCount-len(a.FallObjects); i++ {
			newParticle := FallObject{}
			newParticle.Reset(a.DynamicConfig, scaleX, scaleY)
			newParticle.Position.Y -= float64(i) * a.DynamicConfig.Gravity
			a.FallObjects = append(a.FallObjects, newParticle)
		}
	}
}

func (ob *FallObject) Reset(config *DynamicConfig, scaleX, scaleY float64) {
	ob.Particle.Size = rand.NormFloat64()*config.SizeVariation + config.Size
	if ob.Particle.Size < 1 {
		ob.Particle.Size = 1
	}
	ob.Particle.Position.X = (rand.Float64() * (1000 + (scaleX * 2))) - scaleX
	ob.Particle.Position.Y = -rand.Float64()*100 - ob.Size - scaleY
	ob.Speed = config.Gravity + rand.Float64()*config.Gravity/2
	ob.Particle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	ob.Particle.LineWidth = config.LineWidth
	ob.Particle.Hollow = config.Hollow
	ob.Particle.Rotation = config.Rotation
	ob.Particle.Shape = config.Shape
	ob.Phase = 2 * math.Pi * rand.Float64()
}

func (ob *FallObject) randomSwing(t float64) float64 {
	t += ob.Phase
	amplitudes := []float64{1, 0.5, 0.8}
	frequencies := []float64{2, 3, 1.5}
	phaseShifts := []float64{0, math.Pi / 2, math.Pi}

	result := 0.0
	for i := 0; i < len(amplitudes); i++ {
		result += amplitudes[i] * math.Cos(frequencies[i]*t+phaseShifts[i])
	}
	return result
}

func (a *FallingShape) TriggerWindGust() {
	a.WindDirection = 1
	if rand.Float64() < 0.5 {
		a.WindDirection = -1
	}
	a.WindStrength = rand.Float64() * float64(a.DynamicConfig.WindgutsStrength) * 20
}
