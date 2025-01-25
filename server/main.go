package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// AnimationConfig holds the parameters for the animation
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

var config AnimationConfig

// Mutex-protected map to store WebSocket connections
var clients = make(map[string]*websocket.Conn)
var clientsMutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle incoming WebSocket connection
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Generate a unique client ID
	clientID := r.URL.Query().Get("client_id")
	if clientID == "" {
		clientID = generateUniqueID()
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	// Store the client connection
	clientsMutex.Lock()
	clients[clientID] = conn
	clientsMutex.Unlock()
	log.Printf("Client connected: %s", clientID)

	// Notify the client of its assigned ID
	err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"client_id": "%s"}`, clientID)))
	if err != nil {
		log.Printf("Error sending client ID to %s: %v", clientID, err)
	}

	// Send the current configuration to the newly connected client
	err = sendToClient(clientID, config)
	if err != nil {
		log.Printf("Error sending initial config to %s: %v", clientID, err)
	}

	// Handle client communication
	go handleClientMessages(clientID, conn)
}

// Handle client messages and disconnection
func handleClientMessages(clientID string, conn *websocket.Conn) {
	defer func() {
		clientsMutex.Lock()
		delete(clients, clientID)
		clientsMutex.Unlock()
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

// Handle incoming REST API request to update configuration
func handleUpdateConfig(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON body
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Get the client ID from query parameters
	clientID := r.URL.Query().Get("client_id")

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	if clientID != "" {
		// Send the configuration to the specified client
		err = sendToClient(clientID, config)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending to client %s: %v", clientID, err), http.StatusInternalServerError)
			return
		}
		log.Printf("Configuration sent to client: %s", clientID)
	} else {
		// Broadcast the configuration to all connected clients
		for id := range clients {
			err = sendToClient(id, config)
			if err != nil {
				log.Printf("Error sending to client %s: %v", id, err)
			}
		}
		log.Println("Configuration broadcasted to all clients.")
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Configuration updated successfully.")
}

// Send configuration to a specific client
func sendToClient(clientID string, config AnimationConfig) error {
	conn, ok := clients[clientID]
	if !ok {
		return fmt.Errorf("Client %s is not connected", clientID)
	}

	// Marshal the configuration to JSON
	message, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// Send the configuration to the client
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		clientsMutex.Lock()
		delete(clients, clientID)
		clientsMutex.Unlock()
		return fmt.Errorf("Error sending to client %s: %v", clientID, err)
	}

	return nil
}

// Generate a unique ID for clients
func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func main() {
	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleWebSocket)                               // WebSocket endpoint for clients
	r.HandleFunc("/update-config", handleUpdateConfig).Methods("POST") // REST API endpoint for configuration

	// Serve static files (if you want to serve a web page from this server)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./"))))

	// Start server
	log.Println("Server started at http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
