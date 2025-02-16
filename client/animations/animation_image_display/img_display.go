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
	ImagePath    string  `parameter:"ImagePath,default=resources/image.png"`
	Scale        float32 `parameter:"Scale,default=1.0"`
	ScaleToSound bool    `parameter:"ScaleToSound,default=true"`
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
		return
	}
	var finalScale float32
	if a.DynamicConfig.ScaleToSound {

		// Extract the sound amplitude from the data
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
			soundScale = 0.25 // Prevent the image from becoming too small
		}
		if soundScale > 5.0 {
			soundScale = 5.0 // Prevent excessive scaling

		}

		// Blend dynamic scale with static user-defined scale
		finalScale = a.DynamicConfig.Scale * soundScale

	} else {
		finalScale = a.DynamicConfig.Scale
	}
	// Calculate the new width and height
	width := float32(a.Image.Width) * finalScale
	height := float32(a.Image.Height) * finalScale

	// Position the image at the center of the screen
	x := (500 - width/2)
	y := (500 - height/2)

	// Draw the image to the screen
	rl.DrawTextureEx(a.Image, rl.Vector2{X: x, Y: y}, 0, finalScale, rl.White)
}

// Configure allows updating the configuration for the ImageDisplay (e.g., new image path)
func (a *ImageDisplay) Configure(config preset_animation.AnimationParameters) {
	// Parse new config values into DynamicConfig
	preset_animation.Parse(config, a.DynamicConfig)
	a.update = true

}
