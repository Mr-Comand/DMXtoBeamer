package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpHandler "technikflg.com/dmxtobeamerserver/web"
)

func main() {

	// Setup router
	r := mux.NewRouter()
	// r.HandleFunc("/", httpHandler.ServeIndex).Methods("GET")
	r.HandleFunc("/api/animation/list", httpHandler.HandleAnimationList).Methods("GET")
	r.HandleFunc("/api/animation/get/{animationID}", httpHandler.HandleAnimationGet).Methods("GET")
	r.HandleFunc("/api/animation/image/{animationID}", httpHandler.HandleAnimationGetImage).Methods("GET")
	r.HandleFunc("/api/shader/list", httpHandler.HandleShaderList).Methods("GET")
	r.HandleFunc("/api/shader/get/{shaderID}", httpHandler.HandleShaderGet).Methods("GET")
	r.HandleFunc("/api/shader/image/{shaderID}", httpHandler.HandleShaderGetImage).Methods("GET")
	r.HandleFunc("/api/client/list", httpHandler.HandleClientList).Methods("GET")
	r.HandleFunc("/api/client/set/{clientID}", httpHandler.HandleClientSet).Methods("POST")
	r.HandleFunc("/api/client/set/", httpHandler.HandleClientSet).Methods("POST")
	r.HandleFunc("/api/client/get/{clientID}", httpHandler.HandleClientGet).Methods("GET")

	// WebSocket endpoint
	r.HandleFunc("/ws", httpHandler.HandleWS)

	r.PathPrefix("/").HandlerFunc(httpHandler.StaticFileHandler)

	// Start server
	log.Println("Server started at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", r)
}
