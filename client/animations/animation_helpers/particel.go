package animation_helpers

import (
	"fmt"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape int

const (
	Circle Shape = iota
	Square
	Rectangle
	Triangle
	Snowflake
)

type Particle struct {
	Position  Position
	Color     color.RGBA
	Size      float64
	Shape     Shape
	Hollow    bool
	LineWidth uint8
	Rotation  float64 // Rotation in degrees
}

type Position struct {
	X, Y float64
}

// RotatePoint rotates a point around a given origin
func RotatePoint(x, y, cx, cy, angle float64) (float32, float32) {
	rad := angle * (math.Pi / 180) // Convert degrees to radians
	sin, cos := math.Sin(rad), math.Cos(rad)

	// Translate point to origin, rotate, then translate back
	newX := cos*(x-cx) - sin*(y-cy) + cx
	newY := sin*(x-cx) + cos*(y-cy) + cy

	return float32(newX), float32(newY)
}

func (p *Particle) Draw() {
	if p == nil {
		fmt.Println("Particle not defined")
		return
	}

	switch p.Shape {
	case Circle:
		if p.Hollow {
			for i := int32(0); i < int32(p.LineWidth)+1; i++ {
				rl.DrawCircleLines(int32(p.Position.X), int32(p.Position.Y), float32(p.Size)-float32(i), p.Color)
			}
		} else {
			rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
		}
	case Square, Rectangle:
		width := p.Size
		height := p.Size
		if p.Shape == Rectangle {
			width = p.Size * 1.5
		}

		// Define the four corners of the rectangle/square
		halfW, halfH := width/2, height/2
		p1X, p1Y := RotatePoint(p.Position.X-halfW, p.Position.Y-halfH, p.Position.X, p.Position.Y, p.Rotation)
		p2X, p2Y := RotatePoint(p.Position.X+halfW, p.Position.Y-halfH, p.Position.X, p.Position.Y, p.Rotation)
		p3X, p3Y := RotatePoint(p.Position.X+halfW, p.Position.Y+halfH, p.Position.X, p.Position.Y, p.Rotation)
		p4X, p4Y := RotatePoint(p.Position.X-halfW, p.Position.Y+halfH, p.Position.X, p.Position.Y, p.Rotation)

		if p.Hollow {
			for i := int32(0); i < int32(p.LineWidth)+1; i++ {
				rl.DrawLine(int32(p1X)-i, int32(p1Y)-i, int32(p2X)-i, int32(p2Y)-i, p.Color)
				rl.DrawLine(int32(p2X)-i, int32(p2Y)-i, int32(p3X)-i, int32(p3Y)-i, p.Color)
				rl.DrawLine(int32(p3X)-i, int32(p3Y)-i, int32(p4X)-i, int32(p4Y)-i, p.Color)
				rl.DrawLine(int32(p4X)-i, int32(p4Y)-i, int32(p1X)-i, int32(p1Y)-i, p.Color)
			}
		} else {
			rl.DrawTriangle(
				rl.Vector2{X: float32(p2X), Y: float32(p2Y)},
				rl.Vector2{X: float32(p1X), Y: float32(p1Y)},
				rl.Vector2{X: float32(p3X), Y: float32(p3Y)},
				p.Color)
			rl.DrawTriangle(
				rl.Vector2{X: float32(p4X), Y: float32(p4Y)},
				rl.Vector2{X: float32(p3X), Y: float32(p3Y)},
				rl.Vector2{X: float32(p1X), Y: float32(p1Y)},
				p.Color)
		}
	case Triangle:
		// Define triangle vertices
		p1X, p1Y := RotatePoint(p.Position.X, p.Position.Y-p.Size, p.Position.X, p.Position.Y, p.Rotation)
		p2X, p2Y := RotatePoint(p.Position.X-p.Size, p.Position.Y+p.Size, p.Position.X, p.Position.Y, p.Rotation)
		p3X, p3Y := RotatePoint(p.Position.X+p.Size, p.Position.Y+p.Size, p.Position.X, p.Position.Y, p.Rotation)

		if p.Hollow {
			for i := int32(0); i < int32(p.LineWidth)+1; i++ {
				rl.DrawTriangleLines(
					rl.Vector2{X: float32(p1X) - float32(i), Y: float32(p1Y) - float32(i)},
					rl.Vector2{X: float32(p2X) - float32(i), Y: float32(p2Y) + float32(i)},
					rl.Vector2{X: float32(p3X) + float32(i), Y: float32(p3Y) + float32(i)},
					p.Color,
				)
			}
		} else {
			rl.DrawTriangle(rl.Vector2{X: float32(p1X), Y: float32(p1Y)}, rl.Vector2{X: float32(p2X), Y: float32(p2Y)}, rl.Vector2{X: float32(p3X), Y: float32(p3Y)}, p.Color)
		}
	case Snowflake:
		for i := 0; i < 6; i++ {
			angle := float64(i)*math.Pi/3 + (p.Rotation * math.Pi / 180) // Apply rotation
			x1 := p.Position.X + p.Size*math.Cos(angle)
			y1 := p.Position.Y + p.Size*math.Sin(angle)
			for j := int32(0); j < int32(p.LineWidth); j++ {
				rl.DrawLine(
					int32(p.Position.X)-j, int32(p.Position.Y)-j,
					int32(x1)-j, int32(y1)-j, p.Color,
				)
			}
		}
	default:
		rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
	}
}
