package animations

import (
	"image/color"
	"math"
)

type Ani6 struct {
	Animation
	peakVolume       float64
	VisualValueCount int
	BaseAmplitude    int
	Bandwidth        int
	BaseRadius       float64
	FlourCounter     bool
	dataMap          map[int]int
	DynamicConfig    *Ani4Config
}
type Ani6Config struct {
	RingCount     *uint8
	BaseRadius    *float64
	ColorSegments *uint8
	LineWidth     *uint16
}

func (a *Ani6) GenerateDictionary(start, end int) map[int]int {
	result := make(map[int]int)
	for i := 0; i <= end; i++ {
		result[i] = int(math.Floor(math.Abs(float64(i-end/2)) * 1.2 * (float64(start) / 177)))
	}
	return result
}

func (a *Ani6) NormalizeIntensity(r, g, b, peakVolume float64) (int, int, int) {
	maxIntensity := math.Max(r, math.Max(g, b))
	var normalizedR, normalizedG, normalizedB float64

	if maxIntensity > 127 {
		normalizedR = r / maxIntensity
		normalizedG = g / maxIntensity
		normalizedB = b / maxIntensity
	} else {
		normalizedR = r / peakVolume
		normalizedG = g / peakVolume
		normalizedB = b / peakVolume
	}

	return int(normalizedR * 255), int(normalizedG * 255), int(normalizedB * 255)
}

func (a *Ani6) getValue(id int, values []int) float64 {
	return math.Pow(float64(values[a.dataMap[int(id)]])/255, 2) * float64(a.BaseAmplitude)
}
func (a *Ani6) Render(config *AnimationParameters, data *[]float64) {

	values := make([]int, 0, len(*data))
	for _, v := range *data {
		values = append(values, int(v*10))
	}

	a.peakVolume = 1.0

	particles := make([]Particle, a.VisualValueCount)
	for i := 0; i < a.VisualValueCount; i++ {
		value := a.getValue(i, values)
		a.peakVolume = max(a.peakVolume, value)
		distance := a.BaseRadius + (float64(i) / float64(a.VisualValueCount) * 250)
		angle := float64(i)/float64(a.VisualValueCount)*math.Pi*5 + float64(a.Animation.FrameCount)
		x := math.Cos(angle)*distance + 500
		y := math.Sin(angle)*distance + 500
		particles[i] = Particle{Position: Position{X: x, Y: y}, Color: color.RGBA{R: 255, G: 255, B: 255, A: 255}, Size: value / a.peakVolume * 10}
		particles[i].Draw()
	}

	a.peakVolume -= 1
}

func NewAni6(visualValueCount, baseAmplitude, bandwidth, baseRadius int) *Ani6 {
	// Initialize pointers for DynamicConfig
	ringCount := uint8(10)
	colorSegments := uint8(4)
	lineWidth := uint16(1)
	baseRadiusPointer := float64(baseRadius)

	ani := &Ani6{
		VisualValueCount: visualValueCount,
		BaseAmplitude:    baseAmplitude,
		Bandwidth:        bandwidth,
		BaseRadius:       float64(baseRadius), //TODO: WIDTH
		peakVolume:       1,
		DynamicConfig: &Ani4Config{
			RingCount:     &ringCount,
			ColorSegments: &colorSegments,
			LineWidth:     &lineWidth,
			BaseRadius:    &baseRadiusPointer,
		},
	}

	ani.dataMap = GenerateDictionary(ani.Bandwidth, ani.VisualValueCount)
	return ani
}
