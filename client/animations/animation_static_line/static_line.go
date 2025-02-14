package animation_static_line

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type StaticLine struct {
	preset_animation.Animation
	Rotation      float64
	Y             float64
	DynamicConfig *DynamicConfig
}
type DynamicConfig struct {
	Color                 animation_helpers.Color `parameter:"Color,default=#ff0000"`
	LineWidth             uint8                   `parameter:"LineWidth,default=50"`
	Rotation              float64                 `parameter:"Rotation,default=0"`
	RotationSpeed         float64                 `parameter:"RotationSpeed,default=0"`
	PositionSwingDistance float64                 `parameter:"PositionSwingDistance,default=500"`
	PositionSwingSpeed    float64                 `parameter:"PositionSwingSpeed,default=0"`
	PositionSwingEase     int                     `parameter:"PositionSwingEase,default=0"`
}
type StaticLineGenerator struct {
}

func NewStaticLineGenerator() *StaticLineGenerator {
	return &StaticLineGenerator{}
}
func (g *StaticLineGenerator) Unload() {
}
func (g *StaticLineGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {

	StaticShape := StaticLine{DynamicConfig: &DynamicConfig{}}
	(&StaticShape).Configure(config)
	return &StaticShape
}
func (a *StaticLine) Render(data *[]float64, dt float64) {
	a.Rotation = math.Mod(a.DynamicConfig.RotationSpeed*rl.GetTime()+a.DynamicConfig.Rotation, 360.0)
	switch a.DynamicConfig.PositionSwingEase {
	case 0:
		a.Y = math.Sin(a.DynamicConfig.PositionSwingSpeed*rl.GetTime()*6.28) * a.DynamicConfig.PositionSwingDistance
	case 1:
		a.Y = math.Mod(a.DynamicConfig.PositionSwingSpeed*rl.GetTime()*6.28, a.DynamicConfig.PositionSwingDistance) - (a.DynamicConfig.PositionSwingDistance / 2)
	default:
		a.Y = math.Sin(a.DynamicConfig.PositionSwingSpeed*rl.GetTime()*6.28) * a.DynamicConfig.PositionSwingDistance
	}
	statX, startY := animation_helpers.RotatePoint(-1000, a.Y+500, 500, 500, a.Rotation)
	endX, endY := animation_helpers.RotatePoint(2000, a.Y+500, 500, 500, a.Rotation)
	rl.DrawLineEx(rl.Vector2{X: statX, Y: startY}, rl.Vector2{X: endX, Y: endY}, float32(a.DynamicConfig.LineWidth), a.DynamicConfig.Color.RGBA)
}
func (a *StaticLine) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)

}
