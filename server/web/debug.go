//go:build !production

package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func GetWebDataPath() http.FileSystem {
	return http.Dir("./web/webdata/webpage")
}

// GET /animation/list
func HandleAnimationList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Base directory
	baseDir := "./web/webdata/animations"

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
func HandleAnimationGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animationID, exists := vars["animationID"]
	if !exists {
		HttpError(w, "animationID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the config.json file
	configPath := fmt.Sprintf("./web/webdata/animations/%s/config.json", animationID)

	// Check if the config.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		HttpError(w, fmt.Sprintf("Animation %s not found", animationID), http.StatusNotFound)
		return
	}

	// Read and return the config.json
	configFile, err := os.Open(configPath)
	if err != nil {
		HttpError(w, fmt.Sprintf("Error opening config.json: %v", err), http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, configFile); err != nil {
		HttpError(w, fmt.Sprintf("Error reading config.json: %v", err), http.StatusInternalServerError)
	}
}

// GET /api/animation/image/{shaderID}
func HandleAnimationGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animationID, exists := vars["animationID"]
	if !exists {
		HttpError(w, "animationID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the directory containing the image
	dirPath := fmt.Sprintf("./web/webdata/animations/%s", animationID)

	// Supported image formats
	supportedExtensions := []string{".jpg", ".png", ".gif", ".webp"}
	for _, ext := range supportedExtensions {
		imagePath := fmt.Sprintf("%s/image%s", dirPath, ext)
		if _, err := os.Stat(imagePath); err == nil {
			http.ServeFile(w, r, imagePath)
			return
		}
	}

	HttpError(w, fmt.Sprintf("Image for animation %s not found", animationID), http.StatusNotFound)
}

// GET /shader/list
func HandleShaderList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Base directory
	baseDir := "./web/webdata/shaders"

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
func HandleShaderGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderID, exists := vars["shaderID"]
	if !exists {
		HttpError(w, "shaderID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the config.json file
	configPath := fmt.Sprintf("./web/webdata/shaders/%s/config.json", shaderID)

	// Check if the config.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		HttpError(w, fmt.Sprintf("Shader %s not found", shaderID), http.StatusNotFound)
		return
	}

	// Read and return the config.json
	configFile, err := os.Open(configPath)
	if err != nil {
		HttpError(w, fmt.Sprintf("Error opening config.json: %v", err), http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, configFile); err != nil {
		HttpError(w, fmt.Sprintf("Error reading config.json: %v", err), http.StatusInternalServerError)
	}
}

// GET /api/shader/image/{shaderID}
func HandleShaderGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderID, exists := vars["shaderID"]
	if !exists {
		HttpError(w, "shaderID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the directory containing the image
	dirPath := fmt.Sprintf("./web/webdata/shaders/%s", shaderID)

	// Supported image formats
	supportedExtensions := []string{".webp", ".gif", ".jpg", ".png"}
	for _, ext := range supportedExtensions {
		imagePath := fmt.Sprintf("%s/image%s", dirPath, ext)
		if _, err := os.Stat(imagePath); err == nil {
			http.ServeFile(w, r, imagePath)
			return
		}
	}

	HttpError(w, fmt.Sprintf("Image for shader %s not found", shaderID), http.StatusNotFound)
}

// StaticFileHandler serves static files from the os  filesystem
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the requested file path
	filePath := r.URL.Path

	// If the path is just the root ("/"), serve the default "index.html"
	if filePath == "/" {
		filePath = "/index.html"
	}

	// Open the file from the embedded filesystem
	file, err := os.Open("./web/webdata/webpage" + filePath)
	if err != nil {
		// If the file is not found, return a 404
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	// Read the file content into memory
	fileContent, err := io.ReadAll(file)
	if err != nil {
		// If there was an error reading the file, return an internal server error
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Wrap the file content in a bytes.Reader to implement ReadSeeker
	contentReader := bytes.NewReader(fileContent)

	// Get file info (to determine content type)
	fileInfo, err := file.Stat()
	if err != nil {
		// If there was an error getting file info, return an internal server error
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Serve the file content
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), contentReader)
}
