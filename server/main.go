package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GET /animation/list
func handleAnimationList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Base directory
	baseDir := "./webdata/animations"

	// Map to hold the result
	animationData := make(map[string]interface{})

	// Open the base directory
	files, err := os.ReadDir(baseDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading animations directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Iterate over the directories in the base directory
	for _, file := range files {
		if file.IsDir() {
			configPath := fmt.Sprintf("%s/%s/config.json", baseDir, file.Name())
			configFile, err := os.Open(configPath)
			if err != nil {
				log.Printf("Error reading config.json in %s: %v\n", file.Name(), err)
				continue
			}

			var config interface{}
			err = json.NewDecoder(configFile).Decode(&config)
			configFile.Close()
			if err != nil {
				log.Printf("Error parsing config.json in %s: %v\n", file.Name(), err)
				continue
			}

			// Add the parsed config to the result map
			animationData[file.Name()] = config
		}
	}

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(animationData); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// GET /api/animation/get/{shaderID}
func handleAnimationGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animationID, exists := vars["shaderID"]
	if !exists {
		httpError(w, "animationID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the config.json file
	configPath := fmt.Sprintf("./webdata/animations/%s/config.json", animationID)

	// Check if the config.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		httpError(w, fmt.Sprintf("Animation %s not found", animationID), http.StatusNotFound)
		return
	}

	// Read and return the config.json
	configFile, err := os.Open(configPath)
	if err != nil {
		httpError(w, fmt.Sprintf("Error opening config.json: %v", err), http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, configFile); err != nil {
		httpError(w, fmt.Sprintf("Error reading config.json: %v", err), http.StatusInternalServerError)
	}
}

// GET /api/animation/image/{shaderID}
func handleAnimationGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animationID, exists := vars["shaderID"]
	if !exists {
		httpError(w, "animationID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the directory containing the image
	dirPath := fmt.Sprintf("./webdata/animations/%s", animationID)

	// Supported image formats
	supportedExtensions := []string{".jpg", ".png", ".gif", ".webp"}
	for _, ext := range supportedExtensions {
		imagePath := fmt.Sprintf("%s/image%s", dirPath, ext)
		if _, err := os.Stat(imagePath); err == nil {
			http.ServeFile(w, r, imagePath)
			return
		}
	}

	httpError(w, fmt.Sprintf("Image for animation %s not found", animationID), http.StatusNotFound)
}

// GET /shader/list
func handleShaderList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Base directory
	baseDir := "./webdata/shaders"

	// Map to hold the result
	animationData := make(map[string]interface{})

	// Open the base directory
	files, err := os.ReadDir(baseDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading shaders directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Iterate over the directories in the base directory
	for _, file := range files {
		if file.IsDir() {
			configPath := fmt.Sprintf("%s/%s/config.json", baseDir, file.Name())
			configFile, err := os.Open(configPath)
			if err != nil {
				log.Printf("Error reading config.json in %s: %v\n", file.Name(), err)
				continue
			}

			var config interface{}
			err = json.NewDecoder(configFile).Decode(&config)
			configFile.Close()
			if err != nil {
				log.Printf("Error parsing config.json in %s: %v\n", file.Name(), err)
				continue
			}

			// Add the parsed config to the result map
			animationData[file.Name()] = config
		}
	}

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(animationData); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// GET /api/shader/get/{shaderID}
func handleShaderGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderID, exists := vars["shaderID"]
	if !exists {
		httpError(w, "shaderID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the config.json file
	configPath := fmt.Sprintf("./webdata/shaders/%s/config.json", shaderID)

	// Check if the config.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		httpError(w, fmt.Sprintf("Shader %s not found", shaderID), http.StatusNotFound)
		return
	}

	// Read and return the config.json
	configFile, err := os.Open(configPath)
	if err != nil {
		httpError(w, fmt.Sprintf("Error opening config.json: %v", err), http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, configFile); err != nil {
		httpError(w, fmt.Sprintf("Error reading config.json: %v", err), http.StatusInternalServerError)
	}
}

// GET /api/shader/image/{shaderID}
func handleShaderGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderID, exists := vars["shaderID"]
	if !exists {
		httpError(w, "shaderID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the directory containing the image
	dirPath := fmt.Sprintf("./webdata/shaders/%s", shaderID)

	// Supported image formats
	supportedExtensions := []string{".webp", ".gif", ".jpg", ".png"}
	for _, ext := range supportedExtensions {
		imagePath := fmt.Sprintf("%s/image%s", dirPath, ext)
		if _, err := os.Stat(imagePath); err == nil {
			http.ServeFile(w, r, imagePath)
			return
		}
	}

	httpError(w, fmt.Sprintf("Image for shader %s not found", shaderID), http.StatusNotFound)
}

// POST /client/set/{clientID}
func handleClientSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, exists := vars["clientID"]
	if !exists {
		clientID = ""
	}
	var clientConfig ClientConfig

	if err := json.NewDecoder(r.Body).Decode(&clientConfig); err != nil {
		httpError(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
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
			httpError(w, fmt.Sprintf("Client %s is not connected", clientID), http.StatusBadRequest)
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
func handleClientGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, exists := vars["clientID"]
	if !exists {
		httpError(w, "clientID is required in the URL path", http.StatusBadRequest)
		return
	}

	clientMapMutex.Lock()
	config, exists := clientConfigs[clientID]
	clientMapMutex.Unlock()

	if !exists {

		httpError(w, fmt.Sprintf("Client %s has no Config.", clientID), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}
func httpError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	fmt.Fprintln(w, "{\"error\":\""+error+"\"}")
}
func handleWS(w http.ResponseWriter, r *http.Request) {
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
}

// Serve the index.html file
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./webdata/webpage/index.html")
}

func main() {
	// Check if index.html exists
	if _, err := os.Stat("./webdata/webpage"); os.IsNotExist(err) {
		log.Fatal("webdata/webpage directory not found in the current directory")
	}

	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/", serveIndex).Methods("GET")
	r.HandleFunc("/api/animation/list", handleAnimationList).Methods("GET")
	r.HandleFunc("/api/animation/get/{shaderID}", handleAnimationGet).Methods("GET")
	r.HandleFunc("/api/animation/image/{shaderID}", handleAnimationGetImage).Methods("GET")
	r.HandleFunc("/api/shader/list", handleShaderList).Methods("GET")
	r.HandleFunc("/api/shader/get/{shaderID}", handleShaderGet).Methods("GET")
	r.HandleFunc("/api/shader/image/{shaderID}", handleShaderGetImage).Methods("GET")
	r.HandleFunc("/api/client/list", handleClientList).Methods("GET")
	r.HandleFunc("/api/client/set/{clientID}", handleClientSet).Methods("POST")
	r.HandleFunc("/api/client/set/", handleClientSet).Methods("POST")
	r.HandleFunc("/api/client/get/{clientID}", handleClientGet).Methods("GET")

	// WebSocket endpoint
	r.HandleFunc("/ws", handleWS)

	staticFileHandler := http.StripPrefix("/", http.FileServer(http.Dir("./webdata/webpage")))
	r.PathPrefix("/").Handler(staticFileHandler)

	// Start server
	log.Println("Server started at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", r)
}
