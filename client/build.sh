#!/bin/bash
set -e  # Exit on any error

echo "Building the Go client..."

OUTPUT_DIR="build"
DLL_FILE="./resources/portaudio.dll"
SHADERS="./shaders"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build the Go executable for Windows
go build -tags production -o "$OUTPUT_DIR/client.exe" main.go || { echo "Go build failed"; exit 1; }

# Copy PortAudio DLL
if [ -f "$DLL_FILE" ]; then
    cp "$DLL_FILE" "$OUTPUT_DIR/"
else
    echo "Warning: $DLL_FILE not found. Make sure PortAudio DLL is in resources."
fi

# Copy shaders directory
if [ -d "$SHADERS" ]; then
    mkdir -p "$OUTPUT_DIR/shaders"
    cp -r "$SHADERS/"* "$OUTPUT_DIR/shaders/"
else
    echo "Warning: Shaders directory $SHADERS not found."
fi

# Create zip archive
cd "$OUTPUT_DIR"
zip -r client.zip ./*
cd ..

echo "Build complete. Check $OUTPUT_DIR for client.exe, portaudio.dll, shaders, and client.zip."
