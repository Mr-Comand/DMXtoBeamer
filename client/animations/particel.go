package animations

import (
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Particle struct {
	Position Position
	Color    color.RGBA
	Size     float64
}

type Position struct {
	X, Y float64
}

func asColor(r, g, b int) color.RGBA {
	// Simple color normalization based on max intensity
	peakVolume := math.Max(float64(r), math.Max(float64(g), float64(b)))
	// maxIntensity := math.Max(float64(r), math.Max(float64(g), float64(b)))
	normalizedRed := float64(r) / peakVolume
	normalizedGreen := float64(g) / peakVolume
	normalizedBlue := float64(b) / peakVolume

	return color.RGBA{
		R: uint8(normalizedRed * 255),
		G: uint8(normalizedGreen * 255),
		B: uint8(normalizedBlue * 255),
		A: 255,
	}
}

func (p *Particle) Draw() {
	if p == nil {
		fmt.Errorf("Partikel not defined")
		return
	}
	rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
	// fmt.Println(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
}
