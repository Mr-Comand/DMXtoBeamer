package animations

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape8 struct {
	Animation
}

func NewShape8() *Shape8 {
	ani := &Shape8{Animation{FrameCount: 0}}
	return ani
}

// Animate8Shape handles the 8-shape animation
func (a *Shape8) Render(config *AnimationParameters, audioData *[]float64) {
	centerX := float32(800) / 2
	centerY := float32(600) / 2
	radius := float32(100)
	speed := float32(0.05) // Slow down the speed for smoother motion

	// Update frameCount to animate the movement
	a.FrameCount += speed

	// Draw two circles that form the "8" shape
	for i := 0; i < 2; i++ {
		angle := a.FrameCount + float32(i)*math.Pi // Offset each circle by 180 degrees
		x := centerX + radius*float32(math.Sin(float64(angle)))
		y := centerY + radius*float32(math.Cos(float64(angle)))
		rl.DrawCircle(int32(x), int32(y), 10, rl.Red)
	}

	// Reset the frame count if it gets too large to prevent overflow
	if a.FrameCount > math.Pi*2 {
		a.FrameCount = 0
	}
}
func (a *Shape8) Reset() {
	a.FrameCount = 0
}
