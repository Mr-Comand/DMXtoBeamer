package animations

import (
	"technikflg.com/dmxToProjector/animations/animation_balls"
	"technikflg.com/dmxToProjector/animations/animation_circle"
	"technikflg.com/dmxToProjector/animations/animation_spiral"
	"technikflg.com/dmxToProjector/animations/animation_static_shape"
	"technikflg.com/dmxToProjector/animations/animation_text_display"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

var AnimationGenerators map[string]preset_animation.AnimationGenerator

func InitAnimations() {
	AnimationGenerators = make(map[string]preset_animation.AnimationGenerator)
	// Add animations to the map
	AnimationGenerators["balls"] = animation_balls.NewGeneratorBallsAnimation()
	AnimationGenerators["spiral1"] = animation_spiral.NewSpiralGenerator(0)
	AnimationGenerators["spiral2"] = animation_spiral.NewSpiralGenerator(1)
	AnimationGenerators["spiral3"] = animation_spiral.NewSpiralGenerator(2)
	AnimationGenerators["circle"] = animation_circle.NewGeneratorCircleAnimation()
	AnimationGenerators["text"] = animation_text_display.NewTextDisplayGenerator("")
	AnimationGenerators["static"] = animation_static_shape.NewStaticShapeGenerator()
}
func UnloadAnimations() {
	for _, v := range AnimationGenerators {
		v.Unload()
	}
}
