package animation_flash_shape

import (
	"math/rand"
	"time"

	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type FlashShape struct {
	preset_animation.Animation
	Particle                 animation_helpers.Particle
	DynamicConfig            *DynamicConfig
	lastSpectrum             []float64
	lastBeatTime             time.Time
	flashDuration            time.Duration
	autocorrelationThreshold float64
}

type DynamicConfig struct {
	Color         animation_helpers.Color `parameter:"Color,default=#ff0000"`
	Shape         animation_helpers.Shape `parameter:"Shape,default=0"`
	Hollow        bool                    `parameter:"Hollow,default=false"`
	Size          float64                 `parameter:"Size,default=50"`
	LineWidth     uint8                   `parameter:"LineWidth,default=50"`
	Rotation      float64                 `parameter:"Rotation,default=0"`
	FlashDuration float64                 `parameter:"FlashDuration,default=0.1"` // Duration of the flash in seconds
	BPM           float64                 `parameter:"BPM,default=-1"`            // -1 for audio-based flashing
}

type FlashShapeGenerator struct{}

func NewFlashShapeGenerator() *FlashShapeGenerator {
	return &FlashShapeGenerator{}
}

func (g *FlashShapeGenerator) Unload() {}

func (g *FlashShapeGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	flashShape := FlashShape{DynamicConfig: &DynamicConfig{}}
	flashShape.Particle = animation_helpers.Particle{
		Shape:    flashShape.DynamicConfig.Shape,
		Size:     50,
		Position: animation_helpers.Position{X: 500, Y: 500},
	}
	flashShape.Configure(config)
	return &flashShape
}

func (a *FlashShape) Render(data *[]float64, dt float64) {
	currentTime := time.Now()

	// If BPM < 0, use audio-based flashing with beat detection
	if a.DynamicConfig.BPM < 0 {
		if data == nil || len(*data) == 0 {
			return
		}

		if a.detectBeat(*data) {
			// Move shape to a random position on beat
			a.Particle.Position.X = rand.Float64() * 1000 // Adjust based on canvas width
			a.Particle.Position.Y = rand.Float64() * 1000 // Adjust based on canvas height
			a.lastBeatTime = currentTime
		}

		// Keep the flash visible for the configured duration
		if currentTime.Sub(a.lastBeatTime) < a.flashDuration {
			a.Particle.Draw()
		}
	} else if a.DynamicConfig.BPM == 0 {
		// No flashing control if BPM is 0
		// Add custom behavior if needed
	} else {
		// Manual flashing control using BPM
		if currentTime.Sub(a.lastBeatTime) >= time.Duration((1/a.DynamicConfig.BPM)*float64(time.Minute)) {
			// Move to a random position at the start of the flash cycle
			a.Particle.Position.X = rand.Float64() * 1000 // Adjust based on canvas width
			a.Particle.Position.Y = rand.Float64() * 1000 // Adjust based on canvas height
			a.lastBeatTime = currentTime                  // Update the last beat time to current time
		}

		// Flash the shape if we are within the flash duration window
		if currentTime.Sub(a.lastBeatTime) < a.flashDuration {
			a.Particle.Draw()
		}
	}

}

func (a *FlashShape) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	a.Particle.Color = a.DynamicConfig.Color.RGBA
	a.Particle.Size = a.DynamicConfig.Size
	a.Particle.Shape = a.DynamicConfig.Shape
	a.Particle.LineWidth = a.DynamicConfig.LineWidth
	a.Particle.Hollow = a.DynamicConfig.Hollow
	a.Particle.Rotation = a.DynamicConfig.Rotation
	a.flashDuration = time.Duration(a.DynamicConfig.FlashDuration * float64(time.Second))
}
func (a *FlashShape) detectBeat(spectrum []float64) bool {

	// Ensure that the spectrum data is valid
	if len(spectrum) == 0 {
		return false
	}

	// Calculate the average spectrum value
	var sum float64
	for _, value := range spectrum {
		sum += value
	}
	avgSpectrum := sum / float64(len(spectrum))

	// Check if the change in the average spectrum exceeds the threshold (debounced check)
	beatDetected := false
	if len(a.lastSpectrum) > 0 {
		// Calculate the change in the spectrum (difference in averages)
		var prevSum float64
		for _, value := range a.lastSpectrum {
			prevSum += value
		}
		prevAvgSpectrum := prevSum / float64(len(a.lastSpectrum))

		// Detect a beat if the difference between the current and previous averages is significant
		if avgSpectrum-prevAvgSpectrum > a.autocorrelationThreshold {
			beatDetected = true
		}
	}

	// Store the current spectrum as the last one for future comparison
	a.lastSpectrum = spectrum

	return beatDetected
}
