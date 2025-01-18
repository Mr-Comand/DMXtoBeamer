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
	animationMap map[string]animations.AnimationInterface
)

func InitWindow() {
	// Initialize the window with the specified width and height (not fullscreen)
	rl.InitWindow(WindowWidth, WindowHeight, "WebSocket-Controlled Animation")

	// Set the window to be resizable by default (this is the default behavior)
	// No need for rl.SetWindowResizable as it's enabled by default
	rl.SetWindowMaxSize(100000, 1000000)
	rl.SetWindowMinSize(10, 10)
	// rl.SetWindowState(rl.FlagBorderlessWindowedMode)
	rl.SetWindowState(rl.FlagWindowResizable)
	// Make the window borderless and always on top
	rl.SetWindowState(rl.FlagWindowTopmost)

	// Set the initial background to black
	rl.ClearBackground(rl.Black)

	animationMap = animations.InitAnimations()

	// Initialize the microphone audio processor
	bufferSize := 512 * 8 // Number of samples per buffer
	InitAudioProcessor(bufferSize)
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
	// windowHeight := rl.GetScreenHeight()

	// Calculate the scaling factors for the window
	scaleX := float32(windowWidth) / float32(WindowWidth)
	// scaleY := float32(windowHeight) / float32(WindowHeight)

	// Clear the screen
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black) // Set the background to black

	// Begin 2D mode with scaling
	rl.BeginMode2D(rl.NewCamera2D(rl.Vector2{X: 0, Y: 0}, rl.Vector2{X: 0, Y: 0}, 0.0, scaleX))

	// Call the function to process audio and get frequency data
	frequencies := processAudio()

	// Call the function to handle the animation logic based on the configuration
	handleAnimation(config, &frequencies)

	// End 2D mode and drawing
	rl.EndMode2D()

	// End drawing
	rl.EndDrawing()

	// Delay to control frame rate
	// time.Sleep(100 * time.Millisecond)
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
