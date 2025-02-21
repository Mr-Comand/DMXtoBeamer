#!/bin/bash

echo "Building the Go project..."

OUTPUT_DIR="build"
DLL_FILE="./resources/libportaudio.dll"  # Update to use .dll for Windows
SHADERS="./shaders/"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Remove everything from the output directory
# rm -rf "$OUTPUT_DIR"/*  # Uncomment this if you want to clear the directory before building

# Build the Go executable
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/client.exe" main.go  # Update to create a .exe file

# Copy the shared library file to the output directory
cp "$DLL_FILE" "$OUTPUT_DIR/"

# Copy the shaders directory
mkdir -p "$OUTPUT_DIR/shaders/"
cp -r "$SHADERS" "$OUTPUT_DIR/shaders/"

# Create a zip archive of the build
zip -r "$OUTPUT_DIR/client.zip" "$OUTPUT_DIR"/*

echo "Build complete. Check the $OUTPUT_DIR folder for the executable and shared library."