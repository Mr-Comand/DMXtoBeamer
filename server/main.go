package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
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
	r.HandleFunc("/api/elementShader/list", httpHandler.HandleElementShaderList).Methods("GET")
	r.HandleFunc("/api/elementShader/get/{shaderID}", httpHandler.HandleElementShaderGet).Methods("GET")
	r.HandleFunc("/api/elementShader/image/{shaderID}", httpHandler.HandleElementShaderGetImage).Methods("GET")
	r.HandleFunc("/api/textureShader/list", httpHandler.HandleTextureShaderList).Methods("GET")
	r.HandleFunc("/api/textureShader/get/{shaderID}", httpHandler.HandleTextureShaderGet).Methods("GET")
	r.HandleFunc("/api/textureShader/image/{shaderID}", httpHandler.HandleTextureShaderGetImage).Methods("GET")
	r.HandleFunc("/api/client/list", httpHandler.HandleClientList).Methods("GET")
	r.HandleFunc("/api/client/set/{clientID}", httpHandler.HandleClientSet).Methods("POST")
	r.HandleFunc("/api/client/set/", httpHandler.HandleClientSet).Methods("POST")
	r.HandleFunc("/api/client/get/{clientID}", httpHandler.HandleClientGet).Methods("GET")

	// WebSocket endpoint
	r.HandleFunc("/ws", httpHandler.HandleWS)

	r.PathPrefix("/").HandlerFunc(httpHandler.StaticFileHandler)
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	// Start server
	log.Println("Server started at http://127.0.0.1:8080")
	http.ListenAndServe(":8080", corsHandler(r))
}
