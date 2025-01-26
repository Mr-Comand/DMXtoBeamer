package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Parameter struct {
	Type    string `json:"type"`
	Min     *int   `json:"min,omitempty"`
	Max     *int   `json:"max,omitempty"`
	Columns *int   `json:"columns,omitempty"`
	Rows    *int   `json:"rows,omitempty"`
}

type Animation struct {
	AnimationName string               `json:"animationName"`
	Parameters    map[string]Parameter `json:"ParameterName"`
}
type AnimationConfig map[string]interface{}

type Layer struct {
	AnimationID        string                            `json:"animationID"`
	Parameters         AnimationConfig                   `json:"parameters"`
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
}

type ClientConfig struct {
	Layers   []Layer `json:"layers"`
	Dimmer   uint8   `json:"dimmer"`
	HueShift uint16  `json:"hueshift"`
	Rotate   int16   `json:"rotate"`
	Pan      int16   `json:"pan"`
	Tilt     int16   `json:"tilt"`
	Scale    uint8   `json:"scale"`
}

var clientConfigs = make(map[string]ClientConfig)
var animations = make(map[string]Animation)
var clientMap = make(map[string]*websocket.Conn)
var clientMapMutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Generate unique ID
func generateUniqueID() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	id := ""
	for i := 0; i < 32; i++ {
		id += string(chars[rand.Intn(len(chars))])
	}
	return id
}

// GET /client/list
func handleClientList(w http.ResponseWriter, r *http.Request) {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()

	clients := make([]string, 0, len(clientMap))
	for clientID := range clientMap {
		clients = append(clients, clientID)
	}

	response := map[string]interface{}{
		"clients": clients,
	}
	json.NewEncoder(w).Encode(response)
}

// GET /animation/list
func handleAnimationList(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(animations)
}

// POST /client/set
func handleClientSet(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ClientID string  `json:"clientID"`
		Layers   []Layer `json:"layers"`
		Dimmer   uint8   `json:"dimmer"`
		HueShift uint16  `json:"hueshift"`
		Rotate   int16   `json:"rotate"`
		Pan      int16   `json:"pan"`
		Tilt     int16   `json:"tilt"`
		Scale    uint8   `json:"scale"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	clientConfig := ClientConfig{
		Layers:   request.Layers,
		Dimmer:   request.Dimmer,
		HueShift: request.HueShift,
		Rotate:   request.Rotate,
		Pan:      request.Pan,
		Tilt:     request.Tilt,
		Scale:    request.Scale,
	}

	if request.ClientID == "" {
		// Broadcast to all clients
		clientMapMutex.Lock()
		for clientID, conn := range clientMap {
			clientConfigs[clientID] = clientConfig
			// Send the updated configuration to the client
			if err := conn.WriteJSON(clientConfig); err != nil {
				log.Printf("Error sending config to client %s: %v", clientID, err)
			}
		}
		clientMapMutex.Unlock()
		log.Println("Configuration broadcasted to all clients.")
	} else {
		// Update a specific client
		clientMapMutex.Lock()
		conn, exists := clientMap[request.ClientID]
		if !exists {
			http.Error(w, fmt.Sprintf("Client %s is not connected", request.ClientID), http.StatusBadRequest)
			clientMapMutex.Unlock()
			return
		}
		clientConfigs[request.ClientID] = clientConfig

		// Send the updated configuration to the specific client
		if err := conn.WriteJSON(clientConfig); err != nil {
			log.Printf("Error sending config to client %s: %v", request.ClientID, err)
		}
		clientMapMutex.Unlock()
		log.Printf("Configuration updated for client: %s", request.ClientID)
	}

	json.NewEncoder(w).Encode(map[string]string{"success": "Configuration updated successfully."})
}

// POST /client/get
func handleClientGet(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ClientID string `json:"clientID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	clientMapMutex.Lock()
	config, exists := clientConfigs[request.ClientID]
	clientMapMutex.Unlock()

	if !exists {
		http.Error(w, fmt.Sprintf("Client %s does not exista", request.ClientID), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(config)
}

// Serve the index.html file
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}

func main() {
	// Check if index.html exists
	if _, err := os.Stat("./index.html"); os.IsNotExist(err) {
		log.Fatal("index.html not found in the current directory")
	}

	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/", serveIndex).Methods("GET")
	r.HandleFunc("/client/list", handleClientList).Methods("GET")
	r.HandleFunc("/animation/list", handleAnimationList).Methods("GET")
	r.HandleFunc("/client/set", handleClientSet).Methods("POST")
	r.HandleFunc("/client/get", handleClientGet).Methods("POST")

	// WebSocket endpoint
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.URL.Query().Get("client_id")
		if clientID == "" {
			clientID = generateUniqueID()
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}

		clientMapMutex.Lock()
		clientMap[clientID] = conn
		clientMapMutex.Unlock()
		config, exists := clientConfigs[clientID]
		if exists {
			// Send the updated configuration to the client
			if err := conn.WriteJSON(config); err != nil {
				log.Printf("Error sending config to client %s: %v", clientID, err)
			}
		}

		log.Printf("Client connected: %s", clientID)

		defer func() {
			clientMapMutex.Lock()
			delete(clientMap, clientID)
			clientMapMutex.Unlock()
			conn.Close()
			log.Printf("Client disconnected: %s", clientID)
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading from client %s: %v", clientID, err)
				return
			}
		}
	})

	// Start server
	log.Println("Server started at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", r)
}
