package animation_circle

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type CircleAnimation struct {
	preset_animation.Animation
	peakVolume       float64
	VisualValueCount int
	BaseAmplitude    int
	Bandwidth        int
	FlourCounter     bool
	dataMap          map[int]int
	DynamicConfig    *DynamicConfig
}
type DynamicConfig struct {
	RingCount     int     `parameter:"RingCount,default=1"`
	BaseRadius    float64 `parameter:"BaseRadius,default=100"`
	RingDistance  float64 `parameter:"RingDistance,default=150"`
	ColorSegments uint8   `parameter:"ColorSegments,default=1"`
	LineWidth     uint16  `parameter:"LineWidth,default=1"`
}
type CircleGenerator struct {
	VisualValueCount int
	BaseAmplitude    int
	Bandwidth        int
}

func NewGeneratorCircleAnimation() *CircleGenerator {

	return &CircleGenerator{
		VisualValueCount: 200,
		BaseAmplitude:    400,
		Bandwidth:        300,
	}
}
func (g *CircleGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	// Initialize pointers for DynamicConfig
	ani := &CircleAnimation{
		VisualValueCount: g.VisualValueCount,
		BaseAmplitude:    g.BaseAmplitude,
		Bandwidth:        g.Bandwidth,
		peakVolume:       1,
		DynamicConfig:    &DynamicConfig{},
	}
	preset_animation.Parse(config, ani.DynamicConfig)
	fmt.Println(ani.DynamicConfig)
	ani.dataMap = GenerateDictionary(ani.Bandwidth, ani.VisualValueCount)
	return ani

}

func GenerateDictionary(start, end int) map[int]int {
	result := make(map[int]int)
	for i := 0; i <= end; i++ {
		result[i] = int(math.Floor(math.Abs(float64(i-end/2)) * 1.2 * (float64(start) / 177)))
	}
	return result
}
func (a *CircleAnimation) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	fmt.Println()
	fmt.Println(a.DynamicConfig)
	fmt.Println()
}
func NormalizeIntensity(r, g, b, peakVolume float64) (int, int, int) {
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

func (a *CircleAnimation) DrawSmoothLine(coordinates [][3]float64) {
	for i := 1; i < len(coordinates)-1; i++ {
		// red, green, blue := NormalizeIntensity(coordinates[10][2], coordinates[100][2], coordinates[200][2], a.peakVolume)
		// color := fmt.Sprintf("%d %d rgb(%d, %d, %d)", red, green, blue)
		a.peakVolume = math.Max(a.peakVolume, math.Max(coordinates[10][2], math.Max(coordinates[100][2], coordinates[200][2])))
		pos := fmt.Sprintf("%d\t%d\t%d\t%d", int32(coordinates[i][0]), int32(coordinates[i][1]), int32(coordinates[i+1][0]), int32(coordinates[i+1][1]))
		fmt.Print(pos)
	}
	fmt.Println()
}
func (a *CircleAnimation) getValue(id int, values []int) float64 {
	return math.Pow(math.Log10(math.Pow(float64(values[a.dataMap[int(id)]])/255, 2)*float64(a.BaseAmplitude)+1), 2) * 10
}
func (a *CircleAnimation) Render(data *[]float64) {
	values := make([]int, 0, len(*data))
	for _, v := range *data {
		values = append(values, int(v*10))
	}

	a.peakVolume = 1.0

	for j := 0; j <= a.DynamicConfig.RingCount; j++ { //TODO: FX3
		coordinates := make([][3]float64, a.VisualValueCount*2)
		for i := 0; i < a.VisualValueCount; i++ {
			value := a.getValue(i, values)
			distance := a.DynamicConfig.BaseRadius + (float64(j)+1)*a.DynamicConfig.RingDistance + value
			angle := float64(i) / float64(a.VisualValueCount) * math.Pi
			x := math.Cos(angle)*distance + 500
			y := math.Sin(angle)*distance + 500
			x2 := x
			y2 := -math.Sin(angle)*distance + 500
			coordinates[i] = [3]float64{x, y, value}
			coordinates[(a.VisualValueCount*2)-i-1] = [3]float64{x2, y2, value}
		}
		a.Draw(coordinates)
	}
	a.peakVolume -= 1
}

// Draw handles the actual rendering of the visual elements using Raylib
func (a *CircleAnimation) Draw(coordinates [][3]float64) {
	rl.SetLineWidth(float32(a.DynamicConfig.LineWidth))
	for i := 0; i < len(coordinates)-1; i++ {
		offset := +(len(coordinates) / int(a.DynamicConfig.ColorSegments+1)) * int(i/(len(coordinates)/(int(a.DynamicConfig.ColorSegments))))
		// Normalize intensity of the colors
		red, green, blue := NormalizeIntensity(coordinates[(10 + offset)][2],
			coordinates[(a.VisualValueCount/2/int(a.DynamicConfig.ColorSegments+1))+offset][2],
			coordinates[(a.VisualValueCount-10)/int(a.DynamicConfig.ColorSegments+1)+offset][2], a.peakVolume)

		// Create a color in Raylib's Color struct
		color := rl.NewColor(uint8(red), uint8(green), uint8(blue), 255)

		// Draw a line between the current point and the next point
		rl.DrawLine(int32(coordinates[i][0]), int32(coordinates[i][1]), int32(coordinates[i+1][0]), int32(coordinates[i+1][1]), color)
	}
}
