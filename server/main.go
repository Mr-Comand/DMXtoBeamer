package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

var goClientConn *websocket.Conn

// Handle incoming WebSocket connection
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Store the WebSocket connection to the Go client
	goClientConn = conn
	log.Println("Go client connected via WebSocket.")
	err = sendToGoClient(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending to Go client: %v", err), http.StatusInternalServerError)
		return
	}
	// Loop to keep reading messages from Go client
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from Go client WebSocket:", err)
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

	// Relay the configuration to the Go client via WebSocket
	err = sendToGoClient(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending to Go client: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Configuration updated successfully.")
}

// Send configuration to Go client via WebSocket
func sendToGoClient(config AnimationConfig) error {
	if goClientConn == nil {
		return fmt.Errorf("Go client is not connected")
	}

	// Marshal the configuration to JSON
	message, err := json.Marshal(config)
	if err != nil {
		return err
	}

	// Send the configuration to the Go client
	err = goClientConn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleWebSocket)                               // WebSocket endpoint for Go client
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
