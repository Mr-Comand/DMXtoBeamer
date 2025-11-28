#!/bin/bash
set -e

echo "Building the Go client..."

OUTPUT_DIR="build"
DLL_FILE="./resources/libportaudio.dll"
SHADERS="./shaders"

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build the Go executable with explicit CGO flags
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
CGO_CFLAGS="-I$PWD/resources" \
CGO_LDFLAGS="-L$PWD/resources -lportaudio" \
go build -tags production -o "$OUTPUT_DIR/client.exe" main.go || { echo "Go build failed"; exit 1; }

# Copy PortAudio DLL
if [ -f "$DLL_FILE" ]; then
    cp "$DLL_FILE" "$OUTPUT_DIR/"
else
    echo "Warning: $DLL_FILE not found. Make sure libportaudio.dll is in resources."
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

echo "Build complete. Check $OUTPUT_DIR for client.exe, libportaudio.dll, shaders, and client.zip."
