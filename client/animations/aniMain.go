package animations

type AnimationParameters map[string]interface{}
type AnimationInterface interface {
	Render(config *AnimationParameters, audioData *[]float64)
	Reset()
	Configure(config map[string]interface{})
}

type Animation struct {
	FrameCount float32
}

func (a *Animation) Render(config *AnimationParameters, audioData *[]float64) {
}
func (a *Animation) Reset() {
	a.FrameCount = 0
}
func (a *Animation) Configure(config map[string]interface{}) {
}

var Animations map[string]AnimationInterface

func InitAnimations() {
	Animations = make(map[string]AnimationInterface)
	Animations["s1"] = NewAni10(500, 400, 300, 400, SpiralParams{
		A:         0.10,
		B:         0.4,
		Angle:     11,
		Intensity: 0.18,
	})
	Animations["ani10"] = NewAni10(300, 400, 300, 400, SpiralParams{
		A:         1.20,
		B:         0.76,
		Angle:     2.44,
		Intensity: 0.18,
	})
	Animations["ani4"] = NewAni4(900, 400, 300, 100)
	Animations["ani6"] = NewAni6(100, 400, 300, 10)
	Animations["8"] = NewShape8()
	Animations["square"] = NewSquare()
	Animations["balls"] = NewBallsAnimation(10)
}
