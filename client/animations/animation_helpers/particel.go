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
	LineWidth int32 // New field for line thickness
}

type Position struct {
	X, Y float64
}

func (p *Particle) Draw() {
	if p == nil {
		fmt.Errorf("Particle not defined")
		return
	}

	switch p.Shape {
	case Circle:
		if p.Hollow {
			for i := int32(0); i < p.LineWidth; i++ {
				rl.DrawCircleLines(int32(p.Position.X), int32(p.Position.Y), float32(p.Size)-float32(i), p.Color)
			}
		} else {
			rl.DrawCircle(int32(p.Position.X), int32(p.Position.Y), float32(p.Size), p.Color)
		}
	case Square:
		if p.Hollow {
			for i := int32(0); i < p.LineWidth; i++ {
				rl.DrawRectangleLines(
					int32(p.Position.X)-int32(p.Size/2)-i, int32(p.Position.Y)-int32(p.Size/2)-i,
					int32(p.Size)+2*i, int32(p.Size)+2*i, p.Color,
				)
			}
		} else {
			rl.DrawRectangle(
				int32(p.Position.X)-int32(p.Size/2), int32(p.Position.Y)-int32(p.Size/2),
				int32(p.Size), int32(p.Size), p.Color,
			)
		}
	case Rectangle:
		width := p.Size * 1.5
		height := p.Size
		if p.Hollow {
			for i := int32(0); i < p.LineWidth; i++ {
				rl.DrawRectangleLines(
					int32(p.Position.X)-int32(width/2)-i, int32(p.Position.Y)-int32(height/2)-i,
					int32(width)+2*i, int32(height)+2*i, p.Color,
				)
			}
		} else {
			rl.DrawRectangle(
				int32(p.Position.X)-int32(width/2), int32(p.Position.Y)-int32(height/2),
				int32(width), int32(height), p.Color,
			)
		}
	case Triangle:
		p1 := rl.Vector2{X: float32(p.Position.X), Y: float32(p.Position.Y - p.Size)}
		p2 := rl.Vector2{X: float32(p.Position.X - p.Size), Y: float32(p.Position.Y + p.Size)}
		p3 := rl.Vector2{X: float32(p.Position.X + p.Size), Y: float32(p.Position.Y + p.Size)}

		if p.Hollow {
			for i := int32(0); i < p.LineWidth; i++ {
				rl.DrawTriangleLines(
					rl.Vector2{X: p1.X - float32(i), Y: p1.Y - float32(i)},
					rl.Vector2{X: p2.X - float32(i), Y: p2.Y + float32(i)},
					rl.Vector2{X: p3.X + float32(i), Y: p3.Y + float32(i)},
					p.Color,
				)
			}
		} else {
			rl.DrawTriangle(p1, p2, p3, p.Color)
		}
	case Snowflake:
		for i := 0; i < 6; i++ {
			angle := float64(i) * math.Pi / 3
			x1 := p.Position.X + p.Size*math.Cos(angle)
			y1 := p.Position.Y + p.Size*math.Sin(angle)
			for j := int32(0); j < p.LineWidth; j++ {
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
