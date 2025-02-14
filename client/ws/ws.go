package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"technikflg.com/dmxToProjector/animations"
	"technikflg.com/dmxToProjector/animations/preset_animation"
)

var (
	currentConfig ClientConfig
	configMutex   sync.Mutex
	reconnectWait = 1 * time.Second // Time to wait before attempting to reconnect
)

type ClientConfig struct {
	Dimmer   float32       `json:"dimmer"`
	HueShift float32       `json:"hueShift"`
	Rotate   float32       `json:"rotate"`
	Pan      int16         `json:"pan"`
	Tilt     int16         `json:"tilt"`
	Scale    float32       `json:"scale"`
	Layers   []LayerConfig `json:"layers"`
}

type LayerConfig struct {
	LayerID            uint16                               `json:"layerID"`
	AnimationID        string                               `json:"animationID"`
	Parameters         preset_animation.AnimationParameters `json:"parameters"`
	Enabled            bool                                 `json:"enabled"`
	Dimmer             float32                              `json:"dimmer"`
	HueShift           float32                              `json:"hueShift"`
	Rotate             float32                              `json:"rotate"`
	Pan                int16                                `json:"pan"`
	Tilt               int16                                `json:"tilt"`
	Scale              float32                              `json:"scale"`
	Shader             string                               `json:"shader"`
	ShaderParameters   map[string]interface{}               `json:"shaderParameters"`
	TextureShader      map[string]map[string]interface{}    `json:"textureShaders"`
	TextureShaderOrder []string                             `json:"textureShaderOrder"`
	Animation          preset_animation.AnimationInterface
}

// ConnectToWebSocket establishes a WebSocket connection to the given URL with auto-reconnect
func ConnectToWebSocket(url string) *websocket.Conn {
	var ws *websocket.Conn
	var err error

	for {
		ws, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Failed to connect to WebSocket: %v. Retrying in %s...\n", err, reconnectWait)
			time.Sleep(reconnectWait)
			continue
		}

		log.Printf("Connected to WebSocket server at %s\n", url)
		return ws
	}
}

// ListenWebSocket continuously listens for messages from the WebSocket server with auto-reconnect
func ListenWebSocket(url string) {
	var ws *websocket.Conn

	for {
		ws = ConnectToWebSocket(url)
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Error reading from WebSocket: %v. Reconnecting...\n", err)
				_ = ws.Close() // Ensure the connection is closed before reconnecting
				break          // Exit the inner loop to reconnect
			}

			var config ClientConfig
			err = json.Unmarshal(message, &config)
			if err != nil {
				log.Printf("Error unmarshaling config: %v\n", err)
				continue
			}

			updateConfig(config)
		}
	}
}

// updateConfig updates the global configuration in a thread-safe manner
func updateConfig(config ClientConfig) {
	configMutex.Lock()
	defer configMutex.Unlock()
	currentConfig.Dimmer = config.Dimmer
	currentConfig.HueShift = config.HueShift
	currentConfig.Pan = config.Pan
	currentConfig.Tilt = config.Tilt
	currentConfig.Rotate = config.Rotate
	currentConfig.Scale = config.Scale
	Layers := make([]LayerConfig, len(config.Layers))

	for newPosition, newLayer := range config.Layers {

		var oldLayer *LayerConfig
		for _, l := range currentConfig.Layers {
			if l.LayerID == newLayer.LayerID {
				oldLayer = &l
			}
		}
		if oldLayer != nil && oldLayer.AnimationID == newLayer.AnimationID {
			if oldLayer.Animation != nil {
				oldLayer.Animation.Configure(newLayer.Parameters)
				newLayer.Animation = (*oldLayer).Animation
			}
			Layers[newPosition] = newLayer
		} else {
			animationGenerator := animations.AnimationGenerators[newLayer.AnimationID]
			if oldLayer != nil && oldLayer.Animation != nil {
				(*oldLayer).Animation.Unload()
			}
			if animationGenerator != nil {
				animation := animationGenerator.Create(newLayer.Parameters)
				newLayer.Animation = animation
				newLayer.Animation.Configure(newLayer.Parameters)
				Layers[newPosition] = newLayer

				log.Println(newLayer.AnimationID, newLayer, Layers)
			} else {
				Layers[newPosition] = newLayer
			}
		}
	}
	for _, oldLayer := range currentConfig.Layers {
		var used bool
		for _, newLayer := range config.Layers {
			if oldLayer.LayerID == newLayer.LayerID {
				used = true
			}
		}
		if !used {
			if oldLayer.Animation != nil {
				oldLayer.Animation.Unload()
			}
		}
	}
	currentConfig.Layers = Layers
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
