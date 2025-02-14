package animation_frequency_bar

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type FrequencyBar struct {
	preset_animation.Animation
	Rotation      float64
	dataMap       map[int]int
	BaseAmplitude int
	DynamicConfig *DynamicConfig
}
type DynamicConfig struct {
	Color            animation_helpers.Color `parameter:"Color,default=#ff0000"`
	LineWidth        uint8                   `parameter:"LineWidth,default=5"`
	VisualValueCount uint16                  `parameter:"VisualValueCount,default=200"`
	Bandwidth        uint16                  `parameter:"Bandwidth,default=300"`
	Rotation         float64                 `parameter:"Rotation,default=0"`
	RotationSpeed    float64                 `parameter:"RotationSpeed,default=0"`
}
type FrequencyBarGenerator struct {
}

func NewFrequencyBarGenerator() *FrequencyBarGenerator {
	return &FrequencyBarGenerator{}
}
func (g *FrequencyBarGenerator) Unload() {
}
func (g *FrequencyBarGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {

	StaticShape := FrequencyBar{DynamicConfig: &DynamicConfig{}, BaseAmplitude: 400}
	(&StaticShape).Configure(config)
	return &StaticShape
}
func (a *FrequencyBar) getValue(id int, values []int) float64 {
	return math.Pow(math.Log10(math.Pow(float64(values[a.dataMap[int(id)]])/255, 2)*float64(a.BaseAmplitude)+1), 2) * 10
}
func (a *FrequencyBar) Render(data *[]float64, dt float64) {
	a.Rotation = math.Mod((a.DynamicConfig.RotationSpeed*rl.GetTime())+a.DynamicConfig.Rotation, 360.0)
	values := make([]int, 0, len(*data))
	for _, v := range *data {
		values = append(values, int(v*10))
	}
	for i := uint16(0); i < a.DynamicConfig.VisualValueCount; i++ {
		value := a.getValue(int(i), values)
		statX, startY := animation_helpers.RotatePoint(float64(i)*(1200/float64(a.DynamicConfig.VisualValueCount))-100, 500+value, 500, 500, a.Rotation)
		endX, endY := animation_helpers.RotatePoint(float64(i)*(1200/float64(a.DynamicConfig.VisualValueCount))-100, 500-value, 500, 500, a.Rotation)
		rl.DrawLineEx(rl.Vector2{X: statX, Y: startY}, rl.Vector2{X: endX, Y: endY}, float32(a.DynamicConfig.LineWidth), a.DynamicConfig.Color.RGBA)
	}
}
func (a *FrequencyBar) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	a.dataMap = GenerateDictionary(int(a.DynamicConfig.Bandwidth), int(a.DynamicConfig.VisualValueCount))
}

func GenerateDictionary(start, end int) map[int]int {
	result := make(map[int]int)
	for i := 0; i <= end; i++ {
		result[i] = int(math.Floor(math.Abs(float64(i-end/2)) * 1.2 * (float64(start) / 177)))
	}
	return result
}
