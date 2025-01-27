package httphandler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	LayerID            uint16                            `json:"layerID"`
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
var clientMap = make(map[string]*websocket.Conn)
var clientMapMutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Generate unique ID
func GenerateUniqueID() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	id := ""
	for i := 0; i < 32; i++ {
		id += string(chars[rand.Intn(len(chars))])
	}
	return id
}

// GET /client/list
func HandleClientList(w http.ResponseWriter, r *http.Request) {
	clientMapMutex.Lock()
	defer clientMapMutex.Unlock()

	clients := make([]string, 0, len(clientMap))
	for clientID := range clientMap {
		clients = append(clients, clientID)
	}

	response := map[string]interface{}{
		"clients": clients,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// POST /client/set/{clientID}
func HandleClientSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, exists := vars["clientID"]
	if !exists {
		clientID = ""
	}
	var clientConfig ClientConfig

	if err := json.NewDecoder(r.Body).Decode(&clientConfig); err != nil {
		HttpError(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	if clientID == "" {
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
		conn, exists := clientMap[clientID]
		if !exists {
			HttpError(w, fmt.Sprintf("Client %s is not connected", clientID), http.StatusBadRequest)
			clientMapMutex.Unlock()
			return
		}
		clientConfigs[clientID] = clientConfig

		// Send the updated configuration to the specific client
		if err := conn.WriteJSON(clientConfig); err != nil {
			log.Printf("Error sending config to client %s: %v", clientID, err)
		}
		clientMapMutex.Unlock()
		log.Printf("Configuration updated for client: %s", clientID)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"success": "Configuration updated successfully."})
}

// POST /client/get
func HandleClientGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, exists := vars["clientID"]
	if !exists {
		HttpError(w, "clientID is required in the URL path", http.StatusBadRequest)
		return
	}

	clientMapMutex.Lock()
	config, exists := clientConfigs[clientID]
	clientMapMutex.Unlock()

	if !exists {

		HttpError(w, fmt.Sprintf("Client %s has no Config.", clientID), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
func HttpError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	fmt.Fprintln(w, "{\"error\":\""+error+"\"}")
}
func HandleWS(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		clientID = GenerateUniqueID()
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
}

// Serve the index.html file
func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./webdata/webpage/index.html")
}
