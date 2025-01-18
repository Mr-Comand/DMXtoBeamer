package animations

import (
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/ws"
)

type Ani10 struct {
	Animation
	particles        []Particle
	peakVolume       float32
	Spiral           SpiralParams
	Particles        []Particle
	PeakVolume       float64
	VisualValueCount int
	BaseAmplitude    int
	Bandwidth        int
	BaseRadius       int
	FlourCounter     bool
	dataMap          map[int]int
}
type SpiralParams struct {
	A         float64
	B         float64
	Angle     float64
	Intensity float64
}

type Particle struct {
	Position Position
	Color    color.RGBA
	Size     float64
}

type Position struct {
	X, Y float64
}

func NewAni10(canvasWidth, canvasHeight, visualValueCount, baseAmplitude, bandwidth, baseRadius int) *Ani10 {
	ani := &Ani10{
		VisualValueCount: visualValueCount,
		BaseAmplitude:    baseAmplitude,
		Bandwidth:        bandwidth,
		BaseRadius:       baseRadius,
		PeakVolume:       1,
		Spiral: SpiralParams{
			A:         1.20,
			B:         0.76,
			Angle:     2.44,
			Intensity: 0.18,
		},
	}

	// Initialize particles
	ani.Particles = make([]Particle, 2048)
	for i := 0; i < 2048; i++ {
		ani.Particles[i] = Particle{
			Position: Position{X: 0, Y: 0},
			Color:    color.RGBA{R: 0, G: 255, B: 0},
			Size:     2,
		}
	}
	ani.dataMap = generateDictionary(ani.Bandwidth, ani.VisualValueCount)

	return ani
}
func (a *Ani10) Init() {

}
func (a *Ani10) Reset() {
}

func (a *Ani10) Render(config *ws.AnimationConfig, data *[]float64) {
	// fmt.Println("frame")
	// Extract values from the data map
	values := make([]int, 0, len(*data))
	for _, v := range *data {
		values = append(values, int(v*10))
	}

	// Adjust spiral angle for animation effect
	if a.FlourCounter {
		a.Spiral.Angle += 0.0000004 * 4
		if a.Spiral.Angle >= 2.87 {
			a.FlourCounter = false
		}
	} else {
		a.Spiral.Angle -= 0.0000004 * 4
		if a.Spiral.Angle <= 2.85 {
			a.FlourCounter = true
		}
	}

	// Update the particle positions, sizes, and colors based on data values
	for i := range a.Particles {
		particle := &a.Particles[i]
		value := getValue(i, values, a)
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
		particle.Color = asColor(int(getValue(10%len(a.Particles), values, a)), int(getValue(100%len(a.Particles), values, a)), int(getValue(200%len(a.Particles), values, a)))
		particle.Draw()
	}

	// Update peak volume (used for intensity scaling)
	a.PeakVolume -= 1
}

func getValue(id int, values []int, a *Ani10) float64 {
	// Access data based on `id` and values
	return math.Pow(float64(values[a.dataMap[int(id/(len(a.Particles)/len(a.dataMap)))]])/255, 2) * float64(a.BaseAmplitude)
}

func generateDictionary(start, end int) map[int]int {
	result := make(map[int]int)
	for i := 0; i <= end; i++ {
		result[i] = int(math.Abs(float64(i-(end/2))) * float64(start) / 177) // Custom transformation logic
	}
	return result
}

func asColor(r, g, b int) color.RGBA {
	// Simple color normalization based on max intensity
	peakVolume := math.Max(float64(r), math.Max(float64(g), float64(b)))
	// maxIntensity := math.Max(float64(r), math.Max(float64(g), float64(b)))
	normalizedRed := float64(r) / peakVolume
	normalizedGreen := float64(g) / peakVolume
	normalizedBlue := float64(b) / peakVolume

	return color.RGBA{
		R: uint8(normalizedRed * 255),
		G: uint8(normalizedGreen * 255),
		B: uint8(normalizedBlue * 255),
		A: 255,
	}
}

func (p *Particle) Draw() {
	if p == nil {
		fmt.Errorf("Partikel not defined")
		return
	}
	rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
	// fmt.Println(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
}
