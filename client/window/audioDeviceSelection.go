package window

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gordonklaus/portaudio"
)

var scrollOffset float32

// Function to display audio device selection with scrolling
func displayAudioDeviceSelection() {
	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}
	// filter devices witrh input channels
	var inputDevices []*portaudio.DeviceInfo
	for _, device := range devices {
		if device.MaxInputChannels > 0 {
			inputDevices = append(inputDevices, device)
		}
	}

	const popupX, popupY, popupWidth, popupHeight = 20, 50, 500, 400
	const itemHeight = 30

	// Define the visible area for scrolling
	visibleArea := rl.Rectangle{X: popupX, Y: popupY, Width: popupWidth, Height: popupHeight}

	// Draw the background and border for the popup
	rl.DrawRectangleRec(visibleArea, rl.Gray)
	rl.DrawRectangleLinesEx(visibleArea, 2, rl.Black)

	// Adjust the scroll offset based on mouse wheel movement
	scrollOffset += float32(rl.GetMouseWheelMove()) * 20
	if scrollOffset > 0 {
		scrollOffset = 0 // Prevent scrolling past the top
	}

	// Calculate the maximum scroll offset
	maxScrollOffset := float32(len(inputDevices)*itemHeight - int(popupHeight))
	if scrollOffset < -maxScrollOffset {
		scrollOffset = -maxScrollOffset // Prevent scrolling past the bottom
	}

	// Begin scissor mode to limit drawing to the visible area
	rl.BeginScissorMode(int32(visibleArea.X), int32(visibleArea.Y), int32(visibleArea.Width), int32(visibleArea.Height))

	// Draw each device
	for i, device := range inputDevices {
		yPosition := popupY + scrollOffset + float32(i*itemHeight)

		// Only draw items that are visible in the scissor area
		if yPosition+itemHeight >= popupY && yPosition <= popupY+popupHeight {
			deviceRect := rl.Rectangle{X: popupX + 10, Y: yPosition, Width: popupWidth - 20, Height: itemHeight}
			color := rl.White
			if i == deviceIndex {
				color = rl.Green // Highlight the default device
			}
			if rl.CheckCollisionPointRec(rl.GetMousePosition(), deviceRect) {
				color = rl.Red // Highlight hovered device
				// Handle device selection
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					deviceIndex = i
					InitAudioProcessor(512*8, deviceIndex) // Reinitialize with the selected device
					showPopup = false                      // Close the popup
				}
			}
			rl.DrawRectangleRec(deviceRect, rl.LightGray)
			rl.DrawText(device.Name, int32(deviceRect.X+5), int32(deviceRect.Y+5), 20, color)
		}
	}

	// End scissor mode
	rl.EndScissorMode()

	// Close button
	closeButton := rl.Rectangle{X: popupX + popupWidth - 30, Y: popupY - 30, Width: 30, Height: 30}
	rl.DrawRectangleRec(closeButton, rl.Red)
	rl.DrawText("X", int32(closeButton.X+10), int32(closeButton.Y+5), 20, rl.White)
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), closeButton) {
		showPopup = false // Close the popup
	}
}
