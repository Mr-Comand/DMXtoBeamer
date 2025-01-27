//go:build production

package httphandler

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed webdata/webpage/*
var webpageFS embed.FS

//go:embed webdata/animations/*
var animationsFS embed.FS

//go:embed webdata/shaders/*
var shaderFS embed.FS

// GetWebDataPath will now use the embedded filesystem
func GetWebDataPath() http.FileSystem {
	// Create an http.FS from the embedded FS
	return http.FS(webpageFS)
}

// HandleAnimationList handles the GET request for /animation/list
func HandleAnimationList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Use the embedded filesystem instead of reading from the actual file system
	baseDir := "webdata/animations"

	// Map to hold the result
	animationData := make(map[string]interface{})

	// Open the base directory from the embedded FS
	files, err := animationsFS.ReadDir(baseDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading animations directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Iterate over the directories in the base directory
	for _, file := range files {
		if file.IsDir() {
			configPath := fmt.Sprintf("%s/%s/config.json", baseDir, file.Name())
			configFile, err := animationsFS.Open(configPath)
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

// HandleAnimationGet handles the GET request for /api/animation/get/{animationID}
func HandleAnimationGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animationID, exists := vars["animationID"]
	if !exists {
		HttpError(w, "animationID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the config.json file within the embedded filesystem
	configPath := fmt.Sprintf("webdata/animations/%s/config.json", animationID)

	// Try to open the embedded config.json file
	configFile, err := animationsFS.Open(configPath)
	if err != nil {
		HttpError(w, fmt.Sprintf("Animation %s not found", animationID), http.StatusNotFound)
		return
	}
	defer configFile.Close()

	// Decode the JSON content of the config.json
	var config interface{}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		HttpError(w, fmt.Sprintf("Error parsing config.json: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response header for JSON
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	if err := json.NewEncoder(w).Encode(config); err != nil {
		HttpError(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// HandleAnimationGetImage handles the GET request for /api/animation/image/{shaderID}
func HandleAnimationGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animationID, exists := vars["animationID"]
	if !exists {
		HttpError(w, "animationID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the directory containing the image
	dirPath := fmt.Sprintf("webdata/animations/%s", animationID)

	// Supported image formats
	supportedExtensions := []string{".jpg", ".png", ".gif", ".webp"}
	for _, ext := range supportedExtensions {
		imagePath := fmt.Sprintf("%s/image%s", dirPath, ext)

		// Attempt to open the file in the embedded FS
		file, err := animationsFS.Open(imagePath)
		if err == nil {
			// Get the file's metadata (Stat)
			fileInfo, err := file.Stat()
			if err != nil {
				HttpError(w, fmt.Sprintf("Error retrieving file info for %s: %v", imagePath, err), http.StatusInternalServerError)
				return
			}

			// Serve the file using io.ReadAll and http.ServeContent
			content, err := io.ReadAll(file)
			if err != nil {
				HttpError(w, fmt.Sprintf("Error reading file %s: %v", imagePath, err), http.StatusInternalServerError)
				return
			}

			// Set the correct headers and serve the content
			w.Header().Set("Content-Type", http.DetectContentType(content))
			http.ServeContent(w, r, imagePath, fileInfo.ModTime(), bytes.NewReader(content))
			return
		}
		// If file doesn't exist, continue checking the next extension
	}

	HttpError(w, fmt.Sprintf("Image for animation %s not found", animationID), http.StatusNotFound)
}

// HandleShaderList handles the GET request for /shader/list
func HandleShaderList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Use the embedded filesystem instead of reading from the actual file system
	baseDir := "webdata/shaders"

	// Map to hold the result
	shadersData := make(map[string]interface{})

	// Open the base directory from the embedded FS
	files, err := shaderFS.ReadDir(baseDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading shaders directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Iterate over the directories in the base directory
	for _, file := range files {
		if file.IsDir() {
			configPath := fmt.Sprintf("%s/%s/config.json", baseDir, file.Name())
			configFile, err := shaderFS.Open(configPath)
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
			shadersData[file.Name()] = config
		}
	}

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(shadersData); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// HandleShaderGet handles the GET request for /api/shader/get/{shadersID}
func HandleShaderGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderID, exists := vars["shaderID"]
	if !exists {
		HttpError(w, "shadersID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the config.json file within the embedded filesystem
	configPath := fmt.Sprintf("webdata/shaders/%s/config.json", shaderID)

	// Try to open the embedded config.json file
	configFile, err := shaderFS.Open(configPath)
	if err != nil {
		HttpError(w, fmt.Sprintf("Shader %s not found", shaderID), http.StatusNotFound)
		return
	}
	defer configFile.Close()

	// Decode the JSON content of the config.json
	var config interface{}
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		HttpError(w, fmt.Sprintf("Error parsing config.json: %v", err), http.StatusInternalServerError)
		return
	}

	// Set the response header for JSON
	w.Header().Set("Content-Type", "application/json")
	// Write the JSON response
	if err := json.NewEncoder(w).Encode(config); err != nil {
		HttpError(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

// HandleShaderGetImage handles the GET request for /api/shader/image/{shaderID}
func HandleShaderGetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shaderID, exists := vars["shaderID"]
	if !exists {
		HttpError(w, "shaderID is required in the URL path", http.StatusBadRequest)
		return
	}

	// Path to the directory containing the image
	dirPath := fmt.Sprintf("webdata/shaders/%s", shaderID)

	// Supported image formats
	supportedExtensions := []string{".jpg", ".png", ".gif", ".webp"}
	for _, ext := range supportedExtensions {
		imagePath := fmt.Sprintf("%s/image%s", dirPath, ext)

		// Attempt to open the file in the embedded FS
		file, err := shaderFS.Open(imagePath)
		if err == nil {
			// Get the file's metadata (Stat)
			fileInfo, err := file.Stat()
			if err != nil {
				HttpError(w, fmt.Sprintf("Error retrieving file info for %s: %v", imagePath, err), http.StatusInternalServerError)
				return
			}

			// Serve the file using io.ReadAll and http.ServeContent
			content, err := io.ReadAll(file)
			if err != nil {
				HttpError(w, fmt.Sprintf("Error reading file %s: %v", imagePath, err), http.StatusInternalServerError)
				return
			}

			// Set the correct headers and serve the content
			w.Header().Set("Content-Type", http.DetectContentType(content))
			http.ServeContent(w, r, imagePath, fileInfo.ModTime(), bytes.NewReader(content))
			return
		}
		// If file doesn't exist, continue checking the next extension
	}

	HttpError(w, fmt.Sprintf("Image for shader %s not found", shaderID), http.StatusNotFound)
}

// StaticFileHandler serves static files from the embedded filesystem
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the requested file path
	filePath := r.URL.Path

	// If the path is just the root ("/"), serve the default "index.html"
	if filePath == "/" {
		filePath = "/index.html"
	}

	// Open the file from the embedded filesystem
	file, err := webpageFS.Open("webdata/webpage" + filePath)
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
