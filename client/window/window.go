package window

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations"
	"technikflg.com/dmxToProjector/animations/preset_animation"
	"technikflg.com/dmxToProjector/config"
	"technikflg.com/dmxToProjector/window/shaders"
	"technikflg.com/dmxToProjector/ws"
)

var (
	WindowWidth  = 500
	WindowHeight = 500
)

var (
	deviceIndex      int
	showPopup        bool
	audioDataChannel chan []float64
	audioData        []float64
)
var previousTime time.Time

func InitWindow(WindowName string) {
	// Initialize the window with the specified width and height (not fullscreen)
	rl.InitWindow(int32(WindowWidth), int32(WindowHeight), WindowName+" Animation Client")

	// Set the window to be resizable by default (this is the default behavior)
	rl.SetWindowMaxSize(100000, 1000000)
	rl.SetWindowMinSize(10, 10)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetWindowState(rl.FlagWindowTopmost)
	if config.Cfg.Window.Vsync {
		rl.SetConfigFlags(rl.FlagVsyncHint)
	}
	rl.SetConfigFlags(rl.FlagWindowAlwaysRun)
	rl.SetTargetFPS(config.Cfg.Window.MaxFPS) // Cap to 60 FPS
	// Set the initial background to black
	rl.ClearBackground(rl.Black)

	animations.InitAnimations()

	// Initialize the microphone audio processor
	bufferSize := 512 * 8                      // Number of samples per buffer
	deviceIndex = -1                           // Default device index
	audioDataChannel = make(chan []float64, 1) // Channel with buffer size 1 to hold only the latest data
	InitAudioProcessor(bufferSize, deviceIndex)
	shaders.InitShaders(int32(WindowWidth), int32(WindowHeight))
}

func Render() {
	// Get the current configuration in a thread-safe manner
	config := ws.GetConfig()
	// Handle F11 press to toggle between borderless and windowed mode
	if rl.IsKeyPressed(rl.KeyF11) {
		if rl.IsWindowState(rl.FlagBorderlessWindowedMode) {
			// If the window is borderless, set it to windowed mode
			rl.ClearWindowState(rl.FlagBorderlessWindowedMode)
		} else {
			// If the window is not borderless, set it to borderless windowed mode
			rl.SetWindowState(rl.FlagBorderlessWindowedMode)
		}
	}
	if rl.GetScreenWidth() != WindowWidth || rl.GetScreenHeight() != WindowHeight {
		// Update stored values
		WindowWidth = rl.GetScreenWidth()
		WindowHeight = rl.GetScreenHeight()
		// Call your window resize function here
		onWindowResize()
	}
	// Clear the screen
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	// Draw the button to open the audio device selection popup
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.Rectangle{X: 10, Y: 10, Width: 150, Height: 30}) {
		showPopup = true
	}

	// Show the popup if the button is clicked
	if showPopup {
		displayAudioDeviceSelection()
	}
	// rl.BeginMode2D(rl.NewCamera2D(rl.Vector2{X: float32(windowWidth / 2), Y: float32(windowHeight / 2)}, rl.Vector2{X: 500, Y: 500}, 0.0, min(scaleX, scaleY)))

	// Retrieve the latest audio data from the channel (non-blocking)
	select {
	case audioData = <-audioDataChannel: // Get the most recently processed audio data
		// New audio data is available, update audioData
	default:
		if len(audioData) == 0 {
			audioData = <-audioDataChannel
		}
		// No new data available, use the previous data (audioData remains the same)
	}
	// Call the function to handle the animation logic based on the configuration
	handleAnimation(config, &audioData)

	// End drawing
	rl.EndDrawing()
}

// Function to handle animation based on the configuration
func handleAnimation(config ws.ClientConfig, audioData *[]float64) {
	// Get the current window dimensions for scaling

	// Calculate the scaling factors for the window
	scaleX := float32(WindowWidth) / float32(1000)
	scaleY := float32(WindowHeight) / float32(1000)
	currentTime := time.Now()
	dt := currentTime.Sub(previousTime).Seconds()
	previousTime = currentTime
	for _, l := range config.Layers {
		if !l.Enabled {
			continue
		}
		var animation preset_animation.AnimationInterface = l.Animation
		if animation != nil {
			shaders.SetupTextureShader("general", map[string]interface{}{"HueShift": float32(config.HueShift/360) + float32(l.HueShift/360), "Dimmer": (float32(l.Dimmer) / 100) * (float32(config.Dimmer) / 100)})
			shaders.StartTextureShader("general")
			rl.BeginMode2D(rl.NewCamera2D(rl.Vector2{X: float32(WindowWidth/2) + 150, Y: float32(WindowHeight/2) + 150}, rl.Vector2{X: 500 + float32(l.Pan), Y: 500 + float32(l.Tilt)}, float32(l.Rotate), min(scaleX, scaleY)*(float32(config.Scale)/25*float32(l.Scale)/25)))
			if l.Shader != "" {
				shaders.SetupElementShader(l.Shader, l.ShaderParameters)
				shaders.StartElementShader(l.Shader)
			}
			// Create a new slice to store the modified values
			newAudioData := make([]float64, len(*audioData))

			volume := float64(l.Volume)
			if volume <= 0 {
				volume = 1
			}
			// Multiply each value by 5 and store it in the new slice
			for i, v := range *audioData {
				newAudioData[i] = v * 4 * volume

				// fmt.Println(newAudioData[i], v, l.Volume)
			}

			(animation).Render(&newAudioData, dt)
			if l.Shader != "" {
				shaders.EndElementShader(l.Shader)
			}
			rl.EndMode2D()
			if len(l.TextureShaderOrder) > 0 {
				lastShader := "general"
				for _, name := range l.TextureShaderOrder {
					exists := shaders.SetupTextureShader(name, l.TextureShader[name])
					if !exists {
						continue
					}
					shaders.AnotherTextureShader(lastShader, name)
					lastShader = name
				}
				shaders.EndTextureShader(lastShader)

			} else {
				shaders.EndTextureShader("general")
			}
		} else {
			fmt.Println("Unknown animation type:", l.AnimationID, l.Parameters)
		}
	}
}
func onWindowResize() {
	shaders.OnWindowResize(int32(WindowWidth), int32(WindowHeight))
}
