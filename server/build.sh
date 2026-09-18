#!/bin/bash

echo "Building the Go project..."

OUTPUT_DIR="build"
WEBDATA="./web/webdata"

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Remove everything from the output directory
# rm -rf "$OUTPUT_DIR"/*  # Uncomment this if you want to clear the directory before building

# Build the Go executable
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/server.exe" main.go  # Update to create a .exe file

# Copy the shaders directory
mkdir -p "$OUTPUT_DIR/web/webdata"
cp -r "$WEBDATA" "$OUTPUT_DIR/web/"

# Create a zip archive of the build
zip -r "$OUTPUT_DIR/server.zip" "$OUTPUT_DIR"/*

echo "Build complete. Check the $OUTPUT_DIR folder for the executable and shared library."