package animation_static_shape

import (
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type StaticShape struct {
	preset_animation.Animation
	Particle      animation_helpers.Particle
	DynamicConfig *DynamicConfig
}
type DynamicConfig struct {
	Color     animation_helpers.Color `parameter:"Color,default=#ff0000"`
	Shape     animation_helpers.Shape `parameter:"Shape,default=0"`
	Hollow    bool                    `parameter:"Shape,default=false"`
	Size      float64                 `parameter:"Size,default=50"`
	LineWidth uint8                   `parameter:"LineWidth,default=50"`
	Rotation  float64                 `parameter:"Rotation,default=0"`
}
type StaticShapeGenerator struct {
}

func NewStaticShapeGenerator() *StaticShapeGenerator {
	return &StaticShapeGenerator{}
}
func (g *StaticShapeGenerator) Unload() {
}
func (g *StaticShapeGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {

	StaticShape := StaticShape{DynamicConfig: &DynamicConfig{}}
	StaticShape.Particle = animation_helpers.Particle{Shape: StaticShape.DynamicConfig.Shape, Size: 50, Position: animation_helpers.Position{X: 500, Y: 500}}
	(&StaticShape).Configure(config)
	return &StaticShape
}
func (a *StaticShape) Render(data *[]float64, dt float64) {
	a.Particle.Draw()
}
func (a *StaticShape) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	a.Particle.Color = a.DynamicConfig.Color.RGBA
	a.Particle.Size = a.DynamicConfig.Size
	a.Particle.Shape = a.DynamicConfig.Shape
	a.Particle.LineWidth = a.DynamicConfig.LineWidth
	a.Particle.Hollow = a.DynamicConfig.Hollow
	a.Particle.Rotation = a.DynamicConfig.Rotation
}
