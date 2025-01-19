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
	animations["s1"] = NewAni10(500, 400, 300, 400, SpiralParams{
		A:         0.10,
		B:         0.4,
		Angle:     11,
		Intensity: 0.18,
	})
	animations["ani10"] = NewAni10(500, 400, 300, 400, SpiralParams{

		A:         1.20,
		B:         0.76,
		Angle:     2.44,
		Intensity: 0.18,
	})
	animations["ani4"] = NewAni4(900, 400, 300, 100)
	animations["ani6"] = NewAni6(100, 400, 300, 10)
	animations["8"] = NewShape8()
	animations["square"] = NewSquare()
	return animations
}
