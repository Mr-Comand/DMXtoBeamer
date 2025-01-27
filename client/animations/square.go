package animations

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Square struct {
	Animation
}

func NewSquare() *Square {
	ani := &Square{Animation{FrameCount: 0}}
	return ani
}

// Animate8Shape handles the 8-shape animation
func (a *Square) Render(config *AnimationParameters, audioData *[]float64) {
	centerX := float32(800) / 2
	centerY := float32(600) / 2
	sideLength := float32(200)
	speed := float32(0.05)

	// Update frame count to animate the square
	a.FrameCount += speed

	// Calculate movement along the square path
	t := math.Mod(float64(a.FrameCount), 4) // We have 4 sides to the square

	var x, y float32
	if t < 1 {
		// Top side
		x = centerX + sideLength*float32(t)
		y = centerY - sideLength/2
	} else if t < 2 {
		// Right side
		x = centerX + sideLength
		y = centerY - sideLength/2 + sideLength*float32(t-1)
	} else if t < 3 {
		// Bottom side
		x = centerX + sideLength - sideLength*float32(t-2)
		y = centerY + sideLength/2
	} else {
		// Left side
		x = centerX - sideLength/2
		y = centerY + sideLength/2 - sideLength*float32(t-3)
	}

	// Draw the square at the current position
	rl.DrawRectangle(int32(x), int32(y), int32(sideLength), int32(sideLength), rl.Blue)

	// Reset frame count if it exceeds 4 (one loop around the square)
	if a.FrameCount > 4 {
		a.FrameCount = 0
	}
}
func (a *Square) Reset() {
	a.FrameCount = 0
}
