package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"technikflg.com/dmxToProjector/animations"
)

var (
	currentConfig ClientConfig
	configMutex   sync.Mutex
	reconnectWait = 1 * time.Second // Time to wait before attempting to reconnect
)

type ClientConfig struct {
	Dimmer   uint8         `json:"dimmer"`
	HueShift uint16        `json:"hueShift"`
	Rotate   int16         `json:"rotate"`
	Pan      int16         `json:"pan"`
	Tilt     int16         `json:"tilt"`
	Scale    uint8         `json:"scale"`
	Layers   []LayerConfig `json:"layers"`
}

type LayerConfig struct {
	LayerID            uint16                            `json:"layerID"`
	AnimationID        string                            `json:"animationID"`
	Parameters         animations.AnimationParameters    `json:"parameters"`
	Enabled            bool                              `json:"enabled"`
	Dimmer             uint8                             `json:"dimmer"`
	HueShift           uint16                            `json:"hueShift"`
	Rotate             int16                             `json:"rotate"`
	Pan                int16                             `json:"pan"`
	Tilt               int16                             `json:"tilt"`
	Scale              uint8                             `json:"scale"`
	Shader             string                            `json:"shader"`
	ShaderParameters   map[string]interface{}            `json:"shaderParameters"`
	TextureShader      map[string]map[string]interface{} `json:"textureShaders"`
	TextureShaderOrder []string                          `json:"textureShaderOrder"`
	Animation          animations.AnimationInterface
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
				// oldLayer.Dimmer = newLayer.Dimmer
				// oldLayer.Enabled = newLayer.Enabled
				// oldLayer.HueShift = newLayer.HueShift
				// oldLayer.Pan = newLayer.Pan
				// oldLayer.Tilt = newLayer.Tilt
				// oldLayer.Parameters = newLayer.Parameters
				// oldLayer.Rotate = newLayer.Rotate
				// oldLayer.Scale = newLayer.Scale
				// oldLayer.Shader = newLayer.Shader
				// oldLayer.ShaderParameters = newLayer.ShaderParameters
				// oldLayer.TextureShader = newLayer.TextureShader
				// oldLayer.TextureShaderOrder = newLayer.TextureShaderOrder
				newLayer.Animation = (*oldLayer).Animation
			}
			Layers[newPosition] = newLayer
		} else {
			animation := animations.Animations[newLayer.AnimationID]
			if animation != nil {
				animationClone := Clone(animation)
				newLayer.Animation = animationClone
				newLayer.Animation.Configure(newLayer.Parameters)
				Layers[newPosition] = newLayer

				log.Println(newLayer.AnimationID, newLayer, Layers)
			} else {
				Layers[newPosition] = newLayer
			}
		}
	}
	currentConfig.Layers = Layers
	log.Printf("Updated configuration: %+v\n", config)
}

// Clone creates a deep copy of a value pointed to by an interface
func Clone(a animations.AnimationInterface) animations.AnimationInterface {
	if a == nil {
		return nil
	}

	// Use reflection to get the value and type
	originalValue := reflect.ValueOf(a)
	if originalValue.Kind() != reflect.Ptr {
		log.Println("Clone can only handle pointers to structs")
		return nil
	}

	originalValue = originalValue.Elem()
	if originalValue.Kind() != reflect.Struct {
		log.Println("Clone expects a pointer to a struct")
		return nil
	}

	// Create a new instance of the struct
	cloneValue := reflect.New(originalValue.Type()).Elem()

	// Copy the fields from the original to the clone
	cloneValue.Set(originalValue)

	// Return the new instance as the interface
	return cloneValue.Addr().Interface().(animations.AnimationInterface)
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
