package animation_spiral

import (
	"image/color"
	"math"

	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type AniSpiral struct {
	preset_animation.Animation
	Spiral        SpiralParams
	Particles     []animation_helpers.Particle
	PeakVolume    float64
	BaseAmplitude int
	Bandwidth     int
	FlourCounter  bool
	dataMap       map[int]int
	DynamicConfig *DynamicConfig
	tooUpdate     bool
}
type SpiralParams struct {
	A         float64
	B         float64
	Angle     float64
	Intensity float64
}
type DynamicConfig struct {
	Speed            float64 `parameter:"Speed,default=100"`
	FullBright       bool    `parameter:"FullBright,default=false"`
	VisualValueCount int     `parameter:"VisualValueCount,default=30"`
	ParticleCount    int     `parameter:"ParticleCount,default=512"`
}
type SpiralGenerator struct {
	Spiral        SpiralParams
	BaseAmplitude int
	Bandwidth     int
}

func (g *SpiralGenerator) Unload() {}

func (g SpiralGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	ani := &AniSpiral{
		BaseAmplitude: g.BaseAmplitude,
		Bandwidth:     g.Bandwidth,
		PeakVolume:    1,
		Spiral:        g.Spiral,
		DynamicConfig: &DynamicConfig{},
	}
	ani.Configure(config)

	return ani
}

func NewSpiralGenerator(variant uint8) *SpiralGenerator {
	g := SpiralGenerator{}
	g.BaseAmplitude = 400
	g.Bandwidth = 300
	switch variant {
	case 0:
		g.Spiral = SpiralParams{
			A:         1.20,
			B:         0.76,
			Angle:     2.44,
			Intensity: 0.18,
		}
	case 1:
		g.Spiral = SpiralParams{
			A:         0.10,
			B:         0.4,
			Angle:     11,
			Intensity: 0.18,
		}
	case 2:
		g.Spiral = SpiralParams{
			A:         0.70,
			B:         0.5,
			Angle:     5,
			Intensity: 0.18,
		}
	}
	return &g
}

func (a *AniSpiral) Configure(config preset_animation.AnimationParameters) {
	// Parse dynamic configuration
	preset_animation.Parse(config, a.DynamicConfig)
	a.tooUpdate = true
}

func (a *AniSpiral) Reset() {}
func (a *AniSpiral) Render(data *[]float64, dt float64) {
	if a.tooUpdate {
		a.tooUpdate = false

		// Reinitialize particles based on updated ParticleCount
		a.Particles = make([]animation_helpers.Particle, a.DynamicConfig.ParticleCount)
		for i := 0; i < a.DynamicConfig.ParticleCount; i++ {
			a.Particles[i] = animation_helpers.Particle{
				Position: animation_helpers.Position{X: 0, Y: 0},
				Color:    color.RGBA{R: 0, G: 255, B: 0},
				Size:     2,
			}
		}

		// Re-generate the data map based on the updated VisualValueCount and Bandwidth
		a.dataMap = a.generateDictionary(a.Bandwidth, a.DynamicConfig.VisualValueCount)
	}
	// Extract values from the data map
	values := make([]int, 0, len(*data))
	for _, v := range *data {
		values = append(values, int(v*10))
	}

	// Adjust spiral angle for animation effect
	if a.FlourCounter {
		a.Spiral.Angle += 0.0000004 * 4 * dt * a.DynamicConfig.Speed
		if a.Spiral.Angle >= 2.87 {
			a.FlourCounter = false
		}
	} else {
		a.Spiral.Angle -= 0.0000004 * 4 * dt * a.DynamicConfig.Speed
		if a.Spiral.Angle <= 2.85 {
			a.FlourCounter = true
		}
	}

	// Update the particle positions, sizes, and colors based on data values
	for i, particle := range a.Particles {
		// particle := &a.Particles[i]
		pos := float64(i) * (float64(2048) / float64(a.DynamicConfig.ParticleCount))
		value := a.getValue((pos), values)
		// fmt.Print(", ", value)
		// value = (value * 10) * (value * 10)
		// Calculate positions using an Archimedean spiral with a wavy pattern
		particle.Position.X = (((a.Spiral.A+a.Spiral.B*((a.Spiral.Angle/100)*float64(pos)))*
			math.Cos((a.Spiral.Angle/100)*float64(pos)) +
			math.Sin(float64(pos)/(a.Spiral.Angle/100))*17) + 50) * 10
		particle.Position.Y = (((a.Spiral.A+a.Spiral.B*((a.Spiral.Angle/100)*float64(pos)))*
			math.Sin((a.Spiral.Angle/100)*float64(pos)) +
			math.Cos(float64(pos)/(a.Spiral.Angle/100))*17) + 50) * 10

		// Update size and color based on data
		particle.Size = math.Log(float64(value)/10 + 1)
		if a.DynamicConfig.FullBright {
			particle.Color = animation_helpers.AsFullColor(int(a.getValue(float64(10%len(a.Particles)), values)), int(a.getValue(float64(100%len(a.Particles)), values)), int(a.getValue(float64(200%len(a.Particles)), values)))
		} else {
			particle.Color, a.PeakVolume = animation_helpers.AsDynamicColor(int(a.getValue(float64(10%len(a.Particles)), values)), int(a.getValue(float64(100%len(a.Particles)), values)), int(a.getValue(float64(200%len(a.Particles)), values)), a.PeakVolume)
		}
		particle.Draw()
	}

	// Update peak volume (used for intensity scaling)
	a.PeakVolume -= 50 * dt
}

func (a *AniSpiral) getValue(id float64, values []int) float64 {
	// Prevent divide by zero if len(a.Particles) is 0
	if len(a.Particles) == 0 {
		return 0
	}

	// Prevent division by zero if len(a.dataMap) is zero
	if len(a.dataMap) == 0 {
		return 0
	}

	// Access data based on `id` and values
	particleIndex := int(float64(id) / (float64(len(a.Particles)) / float64(len(a.dataMap))))
	// Ensure we don't access out of bounds in dataMap
	if particleIndex >= len(a.dataMap) {
		particleIndex = len(a.dataMap) - 1
	}
	// Access data based on `id` and values
	return math.Pow(float64(values[a.dataMap[particleIndex]])/255, 2) * float64(a.BaseAmplitude)
}

func (a *AniSpiral) generateDictionary(start, end int) map[int]int {
	result := make(map[int]int)
	for i := 0; i <= end; i++ {
		result[i] = int(math.Abs(float64(i-(end/2))) * float64(start) / 177) // Custom transformation logic
	}
	return result
}
