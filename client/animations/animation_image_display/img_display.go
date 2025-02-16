package animation_image_display

import (
	"fmt"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

type ImageDisplay struct {
	preset_animation.Animation
	DynamicConfig *DynamicConfig
	Image         rl.Texture2D
	update        bool
}

type DynamicConfig struct {
	ImagePath string  `parameter:"ImagePath,default=resources/image.png"`
	Scale     float32 `parameter:"Scale,default=1.0"`
}

type ImageDisplayGenerator struct{}

func NewImageDisplayGenerator() *ImageDisplayGenerator {
	return &ImageDisplayGenerator{}
}
func (g *ImageDisplayGenerator) Unload() {

}
func (a *ImageDisplay) Unload() {
	if a.Image.ID != 0 {
		rl.UnloadTexture(a.Image)
	}
}

// Create generates the ImageDisplay based on the config
func (g *ImageDisplayGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	imageDisplay := ImageDisplay{DynamicConfig: &DynamicConfig{}}
	imageDisplay.Configure(config)
	return &imageDisplay
}

// LoadImage loads the image from the path specified in the DynamicConfig
func (a *ImageDisplay) LoadImage() {
	if a.DynamicConfig.ImagePath != "" {
		imagePath := filepath.Clean(a.DynamicConfig.ImagePath)

		// Check if file exists before loading
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Println("Error: Image file not found at", imagePath)
			return
		}

		// Load image
		a.Image = rl.LoadTexture(imagePath)
	}
}

// Render is used to draw the image onto the screen
func (a *ImageDisplay) Render(data *[]float64, dt float64) {
	if a.update {
		a.update = false
		// Unload previous image if one is already loaded
		if a.Image.ID != 0 {
			rl.UnloadTexture(a.Image)
		}

		// Load new image based on the updated DynamicConfig
		a.LoadImage()
	}
	// Ensure the image is loaded
	if a.Image.ID == 0 {
		// Use a default placeholder image or load a basic texture if loading fails
		// Placeholder: rl.LoadTexture("default_image.png") or handle this case as needed
		return
	}

	// Calculate the scaling based on the DynamicConfig scale value
	scale := a.DynamicConfig.Scale
	width := float32(a.Image.Width) * scale
	height := float32(a.Image.Height) * scale

	// Position the image at the center of the screen
	x := (500 - width/2)
	y := (500 - height/2)

	// Draw the image to the screen
	rl.DrawTextureEx(a.Image, rl.Vector2{X: x, Y: y}, 0, scale, rl.White)
}

// Configure allows updating the configuration for the ImageDisplay (e.g., new image path)
func (a *ImageDisplay) Configure(config preset_animation.AnimationParameters) {
	// Parse new config values into DynamicConfig
	preset_animation.Parse(config, a.DynamicConfig)
	a.update = true

}
