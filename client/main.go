package main

import (
	"log"
	"net/http"

	_ "net/http/pprof" // Import the pprof package

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gordonklaus/portaudio"

	// "github.com/quipo/goprofiler/profiler"
	"technikflg.com/dmxToProjector/window"
	"technikflg.com/dmxToProjector/ws"
)

func main() {

	// pprofConf := profiler.Config{
	// 	Prefix:            "./tmp/myapp.",
	// 	CPU:               true,
	// 	Memory:            true,
	// 	Block:             false,
	// 	Goroutine:         false,
	// 	Mutex:             false,
	// 	Interval:          "15s", // one snapshot every 15 seconds,
	// 	CPUProfileRate:    100,   // collect 100 CPU profiling samples per second
	// 	MemoryProfileRate: 1,     // collect information about all allocations
	// }
	// Prof := profiler.NewProfiler(pprofConf)
	// go Prof.Run()
	// // take one last snapshot and clean resources
	// defer Prof.Stop()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	window.InitWindow()
	// Initialize the window
	defer rl.CloseWindow()
	defer portaudio.Terminate()
	// Connect to the WebSocket server
	wsConn, err := ws.ConnectToWebSocket("ws://127.0.0.1:8080/ws")
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer wsConn.Close()

	// Start the WebSocket listener in a goroutine
	go ws.ListenWebSocket(wsConn)

	// Start animation loop
	// Start animation loop
	for !rl.WindowShouldClose() {
		window.Render()
		// Prof.TakeSnapshot()
	}

}
