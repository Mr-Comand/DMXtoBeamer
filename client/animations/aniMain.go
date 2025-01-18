package animations

import "technikflg.com/dmxToProjector/ws"

type AnimationInterface interface {
	Render(config *ws.AnimationConfig, audioData *[]float64)
	Reset()
}

type Animation struct {
	FrameCount float32
}

func (a *Animation) Render(config *ws.AnimationConfig, audioData *[]float64) {
}
func (a *Animation) Reset() {
	a.FrameCount = 0
}

func InitAnimations() map[string]AnimationInterface {
	animations := make(map[string]AnimationInterface)
	animations["ani10"] = NewAni10(1536, 1536, 500, 400, 300, 400)
	animations["8"] = NewShape8()
	animations["square"] = NewSquare()
	return animations
}
