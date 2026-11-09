#!/bin/bash
set -e

echo "Building the Go client..."

OUTPUT_DIR="build"
HEADER_FILE="./resources/portaudio.h"
DLL_FILE="./resources/libportaudio.dll"
SHADERS="./shaders"

# Check required files
if [ ! -f "$HEADER_FILE" ]; then
    echo "Error: $HEADER_FILE not found. Make sure portaudio.h is in resources/"
    exit 1
fi

if [ ! -f "$DLL_FILE" ]; then
    echo "Warning: $DLL_FILE not found. Make sure libportaudio.dll is in resources/"
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build the Go executable with explicit CGO flags.
# NOTE: No NDI® build tag here — the NDI runtime cannot be loaded on
# Windows (purego.Dlopen is unavailable and the proprietary NDI runtime is
# not distributed for use this way), so Windows builds are compiled
# without NDI® support. Package-form build ('.') so all package files
# (e.g. cliflags.go) are included.
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 \
CGO_CFLAGS="-I$PWD/resources" \
CGO_LDFLAGS="-L$PWD/resources -lportaudio" \
go build -tags production -o "$OUTPUT_DIR/client.exe" .

# Copy PortAudio DLL if it exists
if [ -f "$DLL_FILE" ]; then
    cp "$DLL_FILE" "$OUTPUT_DIR/portaudio_x64.dll"
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
