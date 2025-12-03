//go:build !production

package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	_ "net/http/pprof" // Import the pprof package

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gordonklaus/portaudio"

	"technikflg.com/dmxToProjector/animations"
	"technikflg.com/dmxToProjector/artnet"
	"technikflg.com/dmxToProjector/config"
	"technikflg.com/dmxToProjector/window"
	"technikflg.com/dmxToProjector/ws"
)

var ClientId string

func adjustWebSocketURL(input, clientId string) (string, string, bool) {
	// Define the regex for WebSocket URL validation
	urlRegex := regexp.MustCompile(`^(.+?\:\/\/)?(\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3})(:\d+)(\/[^?\n]*?)?(?:\?(.*?&)?client_id=([^&\n]*)?(&.*)?|\?(.*))?$`)
	matches := urlRegex.FindStringSubmatch(input)
	if len(matches) == 0 {
		return "", "", false // Return empty if the input doesn't match at all
	}

	// If group 1 (protocol) is empty, set it to "ws://"
	if matches[1] == "" {
		matches[1] = "ws://"
	}

	// If group 4 (path) is empty, set it to "/ws"
	if matches[4] == "" {
		matches[4] = "/ws"
	}
	query := ""
	if matches[5] == "" && matches[6] == "" && matches[7] == "" && matches[8] == "" {
		if clientId == "" {
			fmt.Print("Please enter the ClientName: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				clientId = scanner.Text()
			}
		}
		query = "client_id=" + clientId
	} else {
		if matches[8] != "" {
			query = matches[8]
			if clientId == "" {
				fmt.Print("Please enter the ClientName: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					clientId = scanner.Text()
				}
			}
			query += "&client_id=" + clientId

		} else {
			query = matches[5]
			if matches[6] == "" {
				if clientId == "" {
					fmt.Print("Please enter the ClientName: ")
					scanner := bufio.NewScanner(os.Stdin)
					if scanner.Scan() {
						clientId = scanner.Text()
					}
				}
				query += "client_id=" + clientId
			} else {
				query += "client_id=" + matches[6]
			}
			query += matches[7]
		}
	}
	// Reconstruct the full URL
	return matches[1] + matches[2] + matches[3] + matches[4] + "?" + query, clientId, true
}
func main() {
	config.Load("config.yaml")

	// Check if a config is provided as a command-line argument
	var valid bool
	ClientId = ""
	wsURL := ""
	if len(os.Args) == 2 {
		wsURL, ClientId, valid = adjustWebSocketURL(os.Args[1], "")
		if !valid {
			ClientId = wsURL
			// If not provided, prompt the user for input
			fmt.Print("Please enter the WebSocket URL: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				wsURL, ClientId, valid = adjustWebSocketURL(scanner.Text(), ClientId)
			}
		}
	} else if len(os.Args) > 2 {
		wsURL, ClientId, valid = adjustWebSocketURL(os.Args[1], os.Args[2])
		if !valid {
			ClientId = os.Args[1]
			wsURL, ClientId, valid = adjustWebSocketURL(os.Args[2], ClientId)
		}
	} else {
		fmt.Print("Please enter the WebSocket URL: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			wsURL, ClientId, valid = adjustWebSocketURL(scanner.Text(), "")
		}
	}

	// Validate the URL (basic check)
	if wsURL == "" || !valid {
		log.Fatal("WebSocket URL must be provided")
	}

	fmt.Println(wsURL)
	artnet.GetArtnet(config.Cfg.Artnet.Interface, "Proj-"+ClientId)
	// Start the pprof HTTP server for diagnostics
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Initialize the window

	defer end()

	window.InitWindow(ClientId)

	// Start the WebSocket listener in a goroutine
	go ws.ListenWebSocket(wsURL)

	// Start animation loop
	for !rl.WindowShouldClose() {
		window.Render()
	}
}
func end() {
	window.StopAudioChannel <- true
	config := ws.GetConfig()
	for _, v := range config.Layers {
		if v.Animation != nil {
			v.Animation.Unload()
			v.Animation = nil
		}
	}
	animations.UnloadAnimations()
	portaudio.Terminate()
	rl.CloseWindow()
}
