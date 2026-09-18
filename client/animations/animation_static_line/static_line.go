package animation_static_line

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type StaticLine struct {
	preset_animation.Animation
	Rotation      float64
	Y             float64
	DynamicConfig *DynamicConfig
}

type DynamicConfig struct {
	Color                 animation_helpers.Color `parameter:"Color,default=#ff0000"`
	LineWidth             uint8                   `parameter:"LineWidth,default=50"`
	Rotation              float64                 `parameter:"Rotation,default=0"`
	RotationSpeed         float64                 `parameter:"RotationSpeed,default=0"`
	PositionSwingDistance float64                 `parameter:"PositionSwingDistance,default=500"`
	PositionSwingSpeed    float64                 `parameter:"PositionSwingSpeed,default=0"`
	PositionSwingEase     int                     `parameter:"PositionSwingEase,default=0"`
	SegmentCount          int                     `parameter:"SegmentCount,default=1"`       // New parameter for line segmentation
	SegmentShape          animation_helpers.Shape `parameter:"SegmentShape,default=0"`       // New parameter for line segmentation
	SegmentSize           float64                 `parameter:"SegmentSize,default=10"`       // New parameter for line segmentation
	SegmentHollow         bool                    `parameter:"SegmentHollow,default=false"`  // New parameter for line segmentation
	SegmentSplits         int                     `parameter:"SegmentSplits,default=0"`      // New parameter for line segmentation
	SegmentSplitsSpeed    float64                 `parameter:"SegmentSplitsSpeed,default=0"` // New parameter for line segmentation
	LineLength            float64                 `parameter:"LineLength,default=1000"`      // New parameter for line segmentation
}

type StaticLineGenerator struct {
}

func NewStaticLineGenerator() *StaticLineGenerator {
	return &StaticLineGenerator{}
}

func (g *StaticLineGenerator) Unload() {
}

func (g *StaticLineGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	StaticShape := StaticLine{DynamicConfig: &DynamicConfig{}}
	(&StaticShape).Configure(config)
	return &StaticShape
}

func (a *StaticLine) Render(data *[]float64, dt float64) {
	a.Rotation = math.Mod(a.DynamicConfig.RotationSpeed*rl.GetTime()+a.DynamicConfig.Rotation, 360.0)

	// Calculate Y-position swing effect
	switch a.DynamicConfig.PositionSwingEase {
	case 0:
		a.Y = math.Sin(a.DynamicConfig.PositionSwingSpeed*rl.GetTime()*6.28) * a.DynamicConfig.PositionSwingDistance
	case 1:
		a.Y = math.Mod(a.DynamicConfig.PositionSwingSpeed*rl.GetTime()*6.28, a.DynamicConfig.PositionSwingDistance) - (a.DynamicConfig.PositionSwingDistance / 2)
	default:
		a.Y = math.Sin(a.DynamicConfig.PositionSwingSpeed*rl.GetTime()*6.28) * a.DynamicConfig.PositionSwingDistance
	}

	// Define start and end points of the full line
	statX, startY := animation_helpers.RotatePoint(500+a.DynamicConfig.LineLength, a.Y+500, 500, 500, a.Rotation)
	endX, endY := animation_helpers.RotatePoint(500-a.DynamicConfig.LineLength, a.Y+500, 500, 500, a.Rotation)

	// Handle segmented line logic
	segmentCount := a.DynamicConfig.SegmentCount
	if segmentCount <= 1 {
		segmentCount = 1 // Ensure at least one segment
		rl.DrawLineEx(rl.Vector2{X: statX, Y: startY}, rl.Vector2{X: endX, Y: endY}, float32(a.DynamicConfig.LineWidth), a.DynamicConfig.Color.RGBA)
		return
	}

	// Calculate segment step size

	particle := animation_helpers.Particle{
		Color:     a.DynamicConfig.Color.RGBA,
		Size:      a.DynamicConfig.SegmentSize,
		Shape:     a.DynamicConfig.SegmentShape,
		Hollow:    a.DynamicConfig.SegmentHollow,
		LineWidth: a.DynamicConfig.LineWidth,
	}
	if a.DynamicConfig.SegmentSplits == 0 {
		stepX := float64(endX-statX) / float64(segmentCount)
		stepY := float64(endY-startY) / float64(segmentCount)
		for i := 0; i < segmentCount; i++ {
			segStartX := float64(statX) + float64(i)*stepX
			segStartY := float64(startY) + float64(i)*stepY
			particle.Position = animation_helpers.Position{X: float64(segStartX), Y: float64(segStartY)}
			particle.Draw()
		}
	} else {
		stepX := float64(endX-statX) / float64(a.DynamicConfig.SegmentSplits)
		stepY := float64(endY-startY) / float64(a.DynamicConfig.SegmentSplits)
		for i := 0; i <= a.DynamicConfig.SegmentSplits; i++ {
			// Draw the line as multiple segments
			segStartX := float64(statX) + float64(i)*stepX + stepX/2*(1-math.Sin(rl.GetTime()*6.28*a.DynamicConfig.SegmentSplitsSpeed))
			segStartY := float64(startY) + float64(i)*stepY + stepY/2*(1-math.Sin(rl.GetTime()*6.28*a.DynamicConfig.SegmentSplitsSpeed))
			subStepX := stepX / float64(segmentCount) * math.Sin(rl.GetTime()*6.28*a.DynamicConfig.SegmentSplitsSpeed)
			subStepY := stepY / float64(segmentCount) * math.Sin(rl.GetTime()*6.28*a.DynamicConfig.SegmentSplitsSpeed)

			for j := 0; j < segmentCount; j++ {
				segStartX += subStepX
				segStartY += subStepY
				particle.Position = animation_helpers.Position{X: float64(segStartX), Y: float64(segStartY)}
				particle.Draw()
			}
		}
	}

}

func (a *StaticLine) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
}
