package animations

import (
	"technikflg.com/dmxToProjector/animations/animation_balls"
	"technikflg.com/dmxToProjector/animations/animation_circle"
	"technikflg.com/dmxToProjector/animations/animation_spiral"
	"technikflg.com/dmxToProjector/animations/animation_text_display"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

var Animations map[string]preset_animation.AnimationGenerator

func InitAnimations() {
	Animations = make(map[string]preset_animation.AnimationGenerator)
	// Add animations to the map
	Animations["balls"] = animation_balls.NewGeneratorBallsAnimation()
	Animations["spiral1"] = animation_spiral.NewSpiralGenerator(0)
	Animations["spiral2"] = animation_spiral.NewSpiralGenerator(1)
	Animations["spiral3"] = animation_spiral.NewSpiralGenerator(2)
	Animations["circle"] = animation_circle.NewGeneratorCircleAnimation()
	Animations["text"] = animation_text_display.NewTextDisplayGenerator("")
}
