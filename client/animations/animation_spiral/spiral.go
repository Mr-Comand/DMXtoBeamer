package animation_spiral

import (
	"image/color"
	"math"

	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type Ani10 struct {
	preset_animation.Animation
	Spiral           SpiralParams
	Particles        []animation_helpers.Particle
	PeakVolume       float64
	VisualValueCount int
	BaseAmplitude    int
	Bandwidth        int
	BaseRadius       int
	FlourCounter     bool
	dataMap          map[int]int
	DynamicConfig    *DynamicConfig
}
type SpiralParams struct {
	A         float64
	B         float64
	Angle     float64
	Intensity float64
}
type DynamicConfig struct {
	Speed      float64 `parameter:"Speed,default=100"`
	FullBright bool    `parameter:"FullBright,default=false"`
}
type SpiralGenerator struct {
	Spiral           SpiralParams
	VisualValueCount int
	BaseAmplitude    int
	Bandwidth        int
	BaseRadius       int
}

func (g SpiralGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {

	ani := &Ani10{
		VisualValueCount: g.VisualValueCount,
		BaseAmplitude:    g.BaseAmplitude,
		Bandwidth:        g.Bandwidth,
		BaseRadius:       g.BaseRadius,
		PeakVolume:       1,
		Spiral:           g.Spiral,
		DynamicConfig:    &DynamicConfig{},
	}
	preset_animation.Parse(config, ani.DynamicConfig)

	// Initialize particles
	ani.Particles = make([]animation_helpers.Particle, 2048)
	for i := 0; i < 2048; i++ {
		ani.Particles[i] = animation_helpers.Particle{
			Position: animation_helpers.Position{X: 0, Y: 0},
			Color:    color.RGBA{R: 0, G: 255, B: 0},
			Size:     2,
		}
	}
	ani.dataMap = ani.generateDictionary(ani.Bandwidth, ani.VisualValueCount)
	return ani
}

func NewSpiralGenerator(variant uint8) *SpiralGenerator {
	g := SpiralGenerator{}
	g.VisualValueCount = 500
	g.BaseAmplitude = 400
	g.Bandwidth = 300
	g.BaseRadius = 400
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
func (a *Ani10) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
}

func (a *Ani10) Reset() {
}

func (a *Ani10) Render(data *[]float64, dt float64) {
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
	for i := range a.Particles {
		particle := &a.Particles[i]
		value := a.getValue(i, values)
		// fmt.Print(", ", value)
		// value = (value * 10) * (value * 10)
		// Calculate positions using an Archimedean spiral with a wavy pattern
		particle.Position.X = (((a.Spiral.A+a.Spiral.B*((a.Spiral.Angle/100)*float64(i)))*
			math.Cos((a.Spiral.Angle/100)*float64(i)) +
			math.Sin(float64(i)/(a.Spiral.Angle/100))*17) + 50) * 10
		particle.Position.Y = (((a.Spiral.A+a.Spiral.B*((a.Spiral.Angle/100)*float64(i)))*
			math.Sin((a.Spiral.Angle/100)*float64(i)) +
			math.Cos(float64(i)/(a.Spiral.Angle/100))*17) + 50) * 10

		// Update size and color based on data
		particle.Size = math.Log(float64(value)/10 + 1)
		if a.DynamicConfig.FullBright {
			particle.Color = animation_helpers.AsFullColor(int(a.getValue(10%len(a.Particles), values)), int(a.getValue(100%len(a.Particles), values)), int(a.getValue(200%len(a.Particles), values)))
		} else {
			particle.Color, a.PeakVolume = animation_helpers.AsDynamicColor(int(a.getValue(10%len(a.Particles), values)), int(a.getValue(100%len(a.Particles), values)), int(a.getValue(200%len(a.Particles), values)), a.PeakVolume)
		}
		particle.Draw()
	}

	// Update peak volume (used for intensity scaling)
	a.PeakVolume -= 50 * dt
}

func (a *Ani10) getValue(id int, values []int) float64 {
	// Access data based on `id` and values
	return math.Pow(float64(values[a.dataMap[int(id/(len(a.Particles)/len(a.dataMap)))]])/255, 2) * float64(a.BaseAmplitude)
}

func (a *Ani10) generateDictionary(start, end int) map[int]int {
	result := make(map[int]int)
	for i := 0; i <= end; i++ {
		result[i] = int(math.Abs(float64(i-(end/2))) * float64(start) / 177) // Custom transformation logic
	}
	return result
}
