package animation_text_display

import (
	"fmt"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type TextDisplay struct {
	preset_animation.Animation
	DynamicConfig *DynamicConfig
	Font          rl.Font
	update        bool
}

type DynamicConfig struct {
	Color        animation_helpers.Color `parameter:"Color,default=#ff0000"`
	Text         string                  `parameter:"Text,default=Hallo Welt!"`
	FontSize     uint16                  `parameter:"FontSize,default=100"`
	FontPath     string                  `parameter:"FontPath,default=resources/font.ttf"`
	ScaleToSound bool                    `parameter:"ScaleToSound,default=true"`
}

type TextDisplayGenerator struct {
}

func NewTextDisplayGenerator(fontPathString string) *TextDisplayGenerator {
	return &TextDisplayGenerator{}
}

func (g *TextDisplayGenerator) Unload() {
}

func (g *TextDisplayGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	textDisplay := TextDisplay{DynamicConfig: &DynamicConfig{}}
	textDisplay.Configure(config)
	return &textDisplay
}

func (a *TextDisplay) Render(data *[]float64, dt float64) {
	// Update font if necessary
	if a.update && a.DynamicConfig.FontPath != "" {
		fontPath := filepath.Clean(a.DynamicConfig.FontPath)
		if a.Font.Texture.ID != 0 {
			rl.UnloadFont(a.Font)
		}

		// Check if the font exists
		if _, err := os.Stat(fontPath); os.IsNotExist(err) {
			fmt.Println("Error: Font file not found at", fontPath)
			a.Font = rl.GetFontDefault()
		}

		// Load new font
		a.Font = rl.LoadFont(fontPath)
		a.update = false
	}

	// Ensure font is loaded
	if a.Font.Texture.ID == 0 {
		a.Font = rl.GetFontDefault() // Use default font if loading fails
	}

	// Calculate final font size based on sound amplitude
	var finalFontSize float32
	if a.DynamicConfig.ScaleToSound {
		var avgAmplitude float64
		if len(*data) > 0 {
			sum := 0.0
			for _, v := range *data {
				sum += v
			}
			avgAmplitude = sum / float64(len(*data)) // Compute the average amplitude
		}

		// Apply sound-based scaling
		soundScale := float32(1.0 + avgAmplitude*0.5) // Adjust scaling factor as needed
		if soundScale < 0.25 {
			soundScale = 0.25 // Prevent the text from becoming too small
		}
		if soundScale > 5.0 {
			soundScale = 5.0 // Prevent excessive scaling
		}

		// Blend dynamic scale with static user-defined font size
		finalFontSize = float32(a.DynamicConfig.FontSize) * soundScale
	} else {
		finalFontSize = float32(a.DynamicConfig.FontSize)
	}

	// Calculate text width and height based on the final font size
	textWidth := rl.MeasureTextEx(a.Font, a.DynamicConfig.Text, finalFontSize, 2).X
	textHeight := finalFontSize // Approximate height

	// Position the text at the center of the screen
	x := (500 - (textWidth / 2))
	y := (500 - (textHeight / 2))

	// Draw the text on the screen
	rl.DrawTextEx(a.Font, a.DynamicConfig.Text, rl.Vector2{X: x, Y: y}, finalFontSize, 2, a.DynamicConfig.Color.RGBA)
}

// Configure allows updating the configuration for the TextDisplay (e.g., new font path or text)
func (a *TextDisplay) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	a.update = true
}
