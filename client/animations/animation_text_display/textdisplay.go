package animation_text_display

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type TextDisplay struct {
	preset_animation.Animation
	DynamicConfig *DynamicConfig
	Font          rl.Font
}
type DynamicConfig struct {
	Text     string `parameter:"Text,default=Hallo Welt!"`
	FontSize uint16 `parameter:"FontSize,default=100"`
}
type TextDisplayGenerator struct {
	Font rl.Font
}

func NewTextDisplayGenerator(fontPathString string) *TextDisplayGenerator {
	var Font rl.Font
	if fontPathString != "" {
		fontPath := filepath.Clean(fontPathString)
		Font = rl.LoadFont(fontPath)
	}
	return &TextDisplayGenerator{Font: Font}
}

func (g *TextDisplayGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	textDisplay := TextDisplay{DynamicConfig: &DynamicConfig{}}
	preset_animation.Parse(config, textDisplay.DynamicConfig)
	textDisplay.Font = g.Font
	// Load the font

	return &textDisplay
}
func (a *TextDisplay) Render(data *[]float64, dt float64) {
	// Ensure font is loaded
	if a.Font.Texture.ID == 0 {
		a.Font = rl.GetFontDefault() // Use default font if loading fails
	}

	textSize := float32(a.DynamicConfig.FontSize)
	textWidth := rl.MeasureTextEx(a.Font, a.DynamicConfig.Text, textSize, 2).X
	textHeight := textSize // Approximate height

	x := (500 - (textWidth / 2))
	y := (500 - (textHeight / 2))

	rl.DrawTextEx(a.Font, a.DynamicConfig.Text, rl.Vector2{X: x, Y: y}, textSize, 2, rl.Red)

}
func (a *TextDisplay) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
}
