#!/bin/bash
set -e

echo "Building the Go client for Linux (native) or Raspberry Pi (cross-compile)..."

# Usage: ./build-linux.sh [target]
#   native  - build for the machine running the script (default)
#   arm64   - build for Raspberry Pi OS 64-bit (aarch64)
#   arm     - build for Raspberry Pi OS 32-bit (armv7l)
TARGET="${1:-native}"

OUTPUT_DIR="build"
SHADERS="./shaders"

case "$TARGET" in
  native)
    GOOS="$(go env GOHOSTOS)"
    GOARCH="$(go env GOHOSTARCH)"
    ;;
  arm64)
    GOOS=linux
    GOARCH=arm64
    ;;
  arm)
    GOOS=linux
    GOARCH=arm
    GOARM=7
    ;;
  *)
    echo "Unknown target '$TARGET'. Use: native | arm64 | arm"
    exit 1
    ;;
esac

echo "Target: GOOS=$GOOS GOARCH=$GOARCH ${GOARM:+GOARM=$GOARM}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

if [ "$GOARCH" = "amd64" ]; then
    # Native build: link against the system PortAudio (requires portaudio
    # headers/libraries, e.g. 'dnf install portaudio-devel' or
    # 'apt install portaudio19-dev' on Raspberry Pi OS / Debian).
    if ! pkg-config --exists portaudio-2.0; then
        echo "Warning: pkg-config could not find portaudio-2.0, relying on default linker paths."
    fi
    CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH \
    CGO_LDFLAGS="$(pkg-config --libs portaudio-2.0 2>/dev/null || printf '%s' '-lportaudio -lasound -lm -lpthread')" \
    go build -tags "production x11" -o "$OUTPUT_DIR/client" .
else
    # Cross compile: use the matching cross toolchain.
    case "$GOARCH" in
      arm64) CC=aarch64-linux-gnu-gcc ;;
      arm)   CC=arm-linux-gnueabihf-gcc ;;
      *)     CC="$GOARCH-linux-gnu-gcc" ;;
    esac
    echo "Cross compiling with CC=$CC (Fedora: 'dnf install gcc-aarch64-linux-gnu' or 'gcc-arm-linux-gnu', Debian/Ubuntu: 'apt install gcc-aarch64-linux-gnu' or 'gcc-arm-linux-gnueabihf')"
    if ! command -v "$CC" >/dev/null 2>&1; then
        echo "Error: cross toolchain '$CC' not found in PATH."
        exit 1
    fi
    if [ "$GOARCH" = "arm" ]; then
        export GOARM=7
    fi
    CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH CC=$CC \
    CGO_LDFLAGS="-lportaudio -lasound -lm -lpthread" \
    go build -tags "production x11" -o "$OUTPUT_DIR/client" .
fi

# Copy shaders directory
if [ -d "$SHADERS" ]; then
    mkdir -p "$OUTPUT_DIR/shaders"
    cp -r "$SHADERS/"* "$OUTPUT_DIR/shaders/"
else
    echo "Warning: Shaders directory $SHADERS not found."
fi

echo "Build complete. Check $OUTPUT_DIR for the client binary and shaders."
echo "Linux runtime dependencies: libportaudio2, alsa-lib (or pulseaudio-alsa/pipewire-alsa)."
