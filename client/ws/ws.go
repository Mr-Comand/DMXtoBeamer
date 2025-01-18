package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	CurrentConfig AnimationConfig
	ConfigMutex   sync.Mutex
)

type AnimationConfig struct {
	Dimmer    int    `json:"dimmer"`
	HueShift  int    `json:"hueshift"`
	Animation string `json:"animation"`
	Fx1       string `json:"fx1"`
	Fx2       string `json:"fx2"`
	Fx3       string `json:"fx3"`
	Fx4       string `json:"fx4"`
	Pan       int    `json:"pan"`
	Tilt      int    `json:"tilt"`
	Rotate    int    `json:"rotate"`
	Zoom      int    `json:"zoom"`
}

// Connect to the WebSocket server
func ConnectToWebSocket(url string) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to WebSocket: %w", err)
	}
	return ws, nil
}

// Listen to WebSocket messages and update config
func ListenWebSocket(ws *websocket.Conn) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading from WebSocket: %v\n", err)
			return
		}

		var config AnimationConfig
		err = json.Unmarshal(message, &config)
		if err != nil {
			log.Printf("Error unmarshaling config: %v\n", err)
			continue
		}

		updateConfig(config)
	}
}

// Update the global configuration
func updateConfig(config AnimationConfig) {
	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()
	CurrentConfig = config
	fmt.Printf("Updated configuration: %+v\n", config)
}

// Get the current configuration in a thread-safe manner
func GetConfig() AnimationConfig {
	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()
	return CurrentConfig
}
