package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	currentConfig ClientConfig
	configMutex   sync.Mutex
)

type ClientConfig struct {
	Dimmer   int           `json:"dimmer"`
	HueShift int           `json:"hueshift"`
	Rotate   int           `json:"rotate"`
	Layers   []LayerConfig `json:"layers"`
}
type AnimationConfig map[string]interface{}
type LayerConfig struct {
	AnimationID string          `json:"animationID"`
	Parameters  AnimationConfig `json:"parameters"`
	Enabled     bool            `json:"enabled"`
}

// ConnectToWebSocket establishes a WebSocket connection to the given URL
func ConnectToWebSocket(url string) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to WebSocket: %w", err)
	}
	log.Printf("Connected to WebSocket server at %s\n", url)
	return ws, nil
}

// ListenWebSocket continuously listens for messages from the WebSocket server
func ListenWebSocket(ws *websocket.Conn) {
	for {
		log.Printf("bbb")
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading from WebSocket: %v\n", err)
			return
		}

		var config ClientConfig
		err = json.Unmarshal(message, &config)
		if err != nil {
			log.Printf("Error unmarshaling config: %v\n", err)
			continue
		}
		log.Printf("aaaaaaaaaaaaaaa", err)

		updateConfig(config)
	}
}

// updateConfig updates the global configuration in a thread-safe manner
func updateConfig(config ClientConfig) {
	configMutex.Lock()
	defer configMutex.Unlock()
	currentConfig = config
	log.Printf("Updated configuration: %+v\n", config)
}

// GetConfig retrieves the current configuration in a thread-safe manner
func GetConfig() ClientConfig {
	configMutex.Lock()
	defer configMutex.Unlock()
	return currentConfig
}

// SendMessage sends a message to the WebSocket server
func SendMessage(ws *websocket.Conn, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	err = ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	log.Printf("Sent message to WebSocket server: %s\n", string(data))
	return nil
}

// CloseWebSocket closes the WebSocket connection gracefully
func CloseWebSocket(ws *websocket.Conn) error {
	err := ws.Close()
	if err != nil {
		return fmt.Errorf("error closing WebSocket: %w", err)
	}

	log.Println("WebSocket connection closed")
	return nil
}
