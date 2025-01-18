package window

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gordonklaus/portaudio"
	"technikflg.com/dmxToProjector/animations"
	"technikflg.com/dmxToProjector/ws"
)

const (
	WindowWidth  = 1536
	WindowHeight = 1536
)

var (
	animationMap map[string]animations.AnimationInterface
	deviceIndex  int
	showPopup    bool
)

func InitWindow() {
	// Initialize the window with the specified width and height (not fullscreen)
	rl.InitWindow(WindowWidth, WindowHeight, "WebSocket-Controlled Animation")

	// Set the window to be resizable by default (this is the default behavior)
	rl.SetWindowMaxSize(100000, 1000000)
	rl.SetWindowMinSize(10, 10)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetWindowState(rl.FlagWindowTopmost)

	// Set the initial background to black
	rl.ClearBackground(rl.Black)

	animationMap = animations.InitAnimations()

	// Initialize the microphone audio processor
	bufferSize := 512 * 8 // Number of samples per buffer
	deviceIndex = -1      // Default device index
	InitAudioProcessor(bufferSize, deviceIndex)

	// List available audio devices
	// listAudioDevices()
}

func listAudioDevices() {
	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	fmt.Println("Available Audio Devices:")
	for i, device := range devices {
		fmt.Printf("%d: %s\n", i, device.Name)
	}
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
	// Call the function to process audio and get frequency data
	frequencies := processAudio()

	// Call the function to handle the animation logic based on the configuration
	handleAnimation(config, &frequencies)

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
