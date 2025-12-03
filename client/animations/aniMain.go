package animations

import (
	"technikflg.com/dmxToProjector/animations/animation_MH"
	"technikflg.com/dmxToProjector/animations/animation_balls"
	"technikflg.com/dmxToProjector/animations/animation_circle"
	"technikflg.com/dmxToProjector/animations/animation_dynamic_line"
	"technikflg.com/dmxToProjector/animations/animation_falling_shape"
	"technikflg.com/dmxToProjector/animations/animation_flash_shape"
	"technikflg.com/dmxToProjector/animations/animation_frequency_bar"
	"technikflg.com/dmxToProjector/animations/animation_full_color"
	"technikflg.com/dmxToProjector/animations/animation_image_display"
	"technikflg.com/dmxToProjector/animations/animation_spiral"
	"technikflg.com/dmxToProjector/animations/animation_static_line"
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
	AnimationGenerators["fall"] = animation_falling_shape.NewFallingShapeGenerator()
	AnimationGenerators["line"] = animation_static_line.NewStaticLineGenerator()
	AnimationGenerators["bars"] = animation_frequency_bar.NewFrequencyBarGenerator()
	AnimationGenerators["dline"] = animation_dynamic_line.NewDynamicLineGenerator()
	AnimationGenerators["flash"] = animation_flash_shape.NewFlashShapeGenerator()
	AnimationGenerators["img"] = animation_image_display.NewImageDisplayGenerator()
	AnimationGenerators["MH"] = animation_MH.NewMHGenerator()
	AnimationGenerators["fullColor"] = animation_full_color.NewDMXColorGenerator()
}
func UnloadAnimations() {
	for _, v := range AnimationGenerators {
		v.Unload()
	}
}
