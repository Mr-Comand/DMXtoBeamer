package main

import (
	"log"
	"net/http"

	_ "net/http/pprof" // Import the pprof package

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gordonklaus/portaudio"

	"technikflg.com/dmxToProjector/window"
	"technikflg.com/dmxToProjector/ws"
)

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	window.InitWindow()
	// Initialize the window
	defer rl.CloseWindow()
	defer portaudio.Terminate()
	// Connect to the WebSocket server
	wsConn, err := ws.ConnectToWebSocket("ws://127.0.0.1:8080/ws") //?client_id=
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer wsConn.Close()

	// Start the WebSocket listener in a goroutine
	go ws.ListenWebSocket(wsConn)

	// Start animation loop
	for !rl.WindowShouldClose() {
		window.Render()
	}

}
