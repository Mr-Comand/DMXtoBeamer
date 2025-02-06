package animation_helpers

import (
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape int

// Define constants for the enum using iota
const (
	Circle Shape = iota // Starts at 0
	Box
	Triangle
	Snowflake
)

type Particle struct {
	Position Position
	Color    color.RGBA
	Size     float64
	Shape    Shape
}

type Position struct {
	X, Y float64
}

func AsFullColor[T int | float64](r, g, b T) color.RGBA {
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
func AsDynamicColor[T int | float64](r, g, b T, peakVolume float64) (color.RGBA, float64) {
	// Normalize the peak volume to ensure the highest RGB component is used
	peakVolume = math.Max(math.Max(float64(r), math.Max(float64(g), float64(b))), peakVolume)

	// Normalize RGB values based on the peak volume
	normalizedRed := float64(r) / peakVolume
	normalizedGreen := float64(g) / peakVolume
	normalizedBlue := float64(b) / peakVolume

	// Return the RGBA color and updated peak volume
	return color.RGBA{
		R: uint8(normalizedRed * 255),
		G: uint8(normalizedGreen * 255),
		B: uint8(normalizedBlue * 255),
		A: 255,
	}, peakVolume
}
func (p *Particle) Draw() {
	if p == nil {
		fmt.Errorf("Partikel not defined")
		return
	}
	switch p.Shape {
	default:
		rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
	case Circle:
		rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
		// case Box:
		// 	rl.DrawPa(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
	}
}
