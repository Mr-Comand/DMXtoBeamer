package window

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations"
	"technikflg.com/dmxToProjector/ws"
)

const (
	WindowWidth  = 1536
	WindowHeight = 1536
)

var (
	animationMap     map[string]animations.AnimationInterface
	deviceIndex      int
	showPopup        bool
	audioDataChannel chan []float64
	audioData        []float64
)

func InitWindow() {
	// Initialize the window with the specified width and height (not fullscreen)
	rl.InitWindow(WindowWidth, WindowHeight, "WebSocket-Controlled Animation")

	// Set the window to be resizable by default (this is the default behavior)
	rl.SetWindowMaxSize(100000, 1000000)
	rl.SetWindowMinSize(10, 10)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetWindowState(rl.FlagWindowTopmost)
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.SetTargetFPS(0) // Cap to 60 FPS
	// Set the initial background to black
	rl.ClearBackground(rl.Black)

	animationMap = animations.InitAnimations()

	// Initialize the microphone audio processor
	bufferSize := 512 * 8                      // Number of samples per buffer
	deviceIndex = -1                           // Default device index
	audioDataChannel = make(chan []float64, 1) // Channel with buffer size 1 to hold only the latest data

	InitAudioProcessor(bufferSize, deviceIndex)

	// Initialize the audio data channel

	// List available audio devices
	// listAudioDevices()
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
	// Get the current window dimensions for scaling
	windowWidth := rl.GetScreenWidth()
	windowHeight := rl.GetScreenHeight()

	// Calculate the scaling factors for the window
	scaleX := float32(windowWidth) / float32(1000)
	scaleY := float32(windowHeight) / float32(1000)

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
	rl.BeginMode2D(rl.NewCamera2D(rl.Vector2{X: float32(windowWidth / 2), Y: float32(windowHeight / 2)}, rl.Vector2{X: 500, Y: 500}, 0.0, min(scaleX, scaleY)))

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
func handleAnimation(config ws.AnimationConfig, audioData *[]float64) {
	animation := animationMap[config.Animation]
	if animation != nil {
		(animation).Render(&config, audioData)
	} else {
		fmt.Println("Unknown animation type:", config.Animation)
	}
}
