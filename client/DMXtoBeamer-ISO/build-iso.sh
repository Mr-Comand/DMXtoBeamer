
#!/usr/bin/env bash
set -Eeuo pipefail

# ============================================================
# Project configuration
# ============================================================

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd -- "$SCRIPT_DIR/.." && pwd)"
ISO_DIR_NAME="$(basename "$SCRIPT_DIR")"

BUILD_SCRIPT="$PROJECT_DIR/build-linux.sh"
BUILD_DIR="$PROJECT_DIR/build"

APP_BINARY="$BUILD_DIR/client"
SHADER_DIR="$BUILD_DIR/shaders"

WORK_DIR="$SCRIPT_DIR/live-build"
ISO_OUTPUT="$SCRIPT_DIR/dmxtobeamer-live.iso"

INSIDE_CONTAINER="${1:-}"

# ============================================================
# Build inside Debian 13 container (works with Fedora host or CI)
# ============================================================

if [[ "$INSIDE_CONTAINER" != "--inside-container" ]]; then

    # Auto-detect container runtime environment
    IS_GITHUB_CI="${GITHUB_ACTIONS:-}"
    
    if [[ -f /etc/fedora-release ]] || [[ -n "$IS_GITHUB_CI" ]]; then
        echo "Detected container-compatible environment (Fedora or GitHub CI)."
        echo "Starting Debian 13 build container."

        # Use Podman or Docker as container runtime
        if ! command -v podman >/dev/null 2>&1; then
            if command -v docker >/dev/null 2>&1; then
                echo "Using Docker as Podman not available"
                CONTAINER_RUNTIME="docker"
            else
                echo "Installing Podman..."
                sudo dnf install -y podman || {
                    echo "Podman/Docker installation failed, trying fallback..."
                    sudo apt-get update && sudo apt-get install -y podman || exit 1
                }
                CONTAINER_RUNTIME="podman"
            fi
        else
            CONTAINER_RUNTIME="podman"
        fi

        $CONTAINER_RUNTIME run --rm --privileged --network host \
            --security-opt label=disable \
            -v "$PROJECT_DIR:/workspace:Z" \
            -w /workspace \
            -e ISO_DIR_NAME="$ISO_DIR_NAME" \
            -e DEBIAN_FRONTEND="noninteractive" \
            docker.io/library/debian:13 \
            bash -c '
                set -Eeuo pipefail
                apt-get update -qq

                apt-get install -y --no-install-recommends \
                    bash \
                    golang-go \
                    gcc \
                    g++ \
                    make \
                    pkg-config \
                    libasound2-dev \
                    portaudio19-dev \
                    libusb-1.0-0-dev \
                    debootstrap \
                    live-build \
                    libudev-dev \
                    libx11-dev \
                    libxrandr-dev \
                    libxcursor-dev \
                    libxinerama-dev \
                    libxi-dev \
                    libxext-dev \
                    libxfixes-dev \
                    libxxf86vm-dev \
                    squashfs-tools \
                    xorriso \
                    dosfstools \
                    parted \
                    rsync \
                    git \
                    ca-certificates \
                    curl \
                    file \
                    binutils \
                    libxkbcommon-dev \
                    libgl-dev

                cd "/workspace/$ISO_DIR_NAME"
                bash ./build-iso.sh --inside-container
            '

        exit $?
    fi
fi

# ============================================================
# Root check
# ============================================================

if [[ "$EUID" -ne 0 ]]; then
    echo "Run this script with:"
    echo "sudo bash build-iso.sh"
    exit 1
fi

# ============================================================
# Validate build environment
# ============================================================

echo
echo "Project directory: $PROJECT_DIR"
echo "ISO directory:     $SCRIPT_DIR"

for command_name in go lb pkg-config; do
    if ! command -v "$command_name" >/dev/null 2>&1; then
        echo "ERROR: Missing command: $command_name"
        exit 1
    fi
done

if [[ ! -f "$BUILD_SCRIPT" ]]; then
    echo "ERROR: Build script not found: $BUILD_SCRIPT"
    exit 1
fi

# ============================================================
# Build Go application
# ============================================================

echo
echo "=========================================="
echo " Building DMXtoBeamer"
echo "=========================================="

cd "$PROJECT_DIR"

bash "$BUILD_SCRIPT" native

if [[ ! -f "$APP_BINARY" ]]; then
    echo "ERROR: Application binary not found: $APP_BINARY"
    exit 1
fi

if [[ ! -d "$SHADER_DIR" ]]; then
    echo "ERROR: Shader directory not found: $SHADER_DIR"
    exit 1
fi

chmod +x "$APP_BINARY"

echo "Application binary:"
file "$APP_BINARY"

# ============================================================
# Prepare Debian Live build
# ============================================================

echo
echo "=========================================="
echo " Preparing Debian 13 Live"
echo "=========================================="

rm -rf "$WORK_DIR"
mkdir -p "$WORK_DIR"

cd "$WORK_DIR"

IMAGE_VARIANT="bookworm"  # Debian 12 (Stable)
lb config \
    --distribution "$IMAGE_VARIANT" \
    --architectures amd64 \
    --binary-images iso-hybrid \
    --archive-areas "main contrib non-free non-free-firmware" \
    --debian-installer none \
    --apt-recommends false \
    --bootappend-live "boot=live components quiet username=kiosk hostname=dmxtobeamer locales=de_DE.UTF-8 keyboard-layouts=de keyboard-model=pc105" \
    --iso-volume "DMXTOBEAMER"

# ============================================================
# Runtime packages
# ============================================================

mkdir -p config/package-lists

cat > config/package-lists/dmxtobeamer.list.chroot <<'EOF'
# Base system
systemd-sysv
sudo
locales
keyboard-configuration
console-setup
dbus
dbus-user-session

# Wayland compositor and seat management
cage
seatd
libseat1
libwayland-client0
libwayland-cursor0
libwayland-egl1
libwayland-server0

# X11 compatibility for the current Raylib build
xwayland
libx11-6
libxrandr2
libxcursor1
libxinerama1
libxi6
libxext6
libxfixes3
libxxf86vm1

# Graphics
libdrm2
libgbm1
libegl1
libgles2
libgl1
mesa-utils

# Networking
network-manager
wpasupplicant
wireless-regdb
iproute2
iputils-ping
curl
ca-certificates

# USB and serial
libusb-1.0-0
libudev1
udev
usbutils
pciutils
libserialport0

# Audio
alsa-utils
alsa-ucm-conf
alsa-topology-conf
libasound2
libportaudio2
pipewire
pipewire-pulse
wireplumber
libspa-0.2-modules

# Firmware
firmware-linux
firmware-linux-nonfree
firmware-misc-nonfree

# Fonts
fonts-dejavu-core

# Editors
nano
vim-tiny

EOF

# ============================================================
# Install application
# ============================================================

APP_DEST="config/includes.chroot/opt/dmxtobeamer"

mkdir -p "$APP_DEST"

install -m 0755 "$APP_BINARY" "$APP_DEST/client"

cp -a "$SHADER_DIR" "$APP_DEST/shaders"

# ============================================================
# System configuration directories
# ============================================================

SYSTEMD_DIR="config/includes.chroot/etc/systemd/system"
BIN_DIR="config/includes.chroot/usr/local/bin"
CONFIG_DIR="config/includes.chroot/etc/dmxtobeamer"

mkdir -p "$SYSTEMD_DIR"
mkdir -p "$BIN_DIR"
mkdir -p "$CONFIG_DIR"

# ============================================================
# Default application configuration
# ============================================================

cat > "$CONFIG_DIR/dmxtobeamer.conf" <<'EOF'
#!/bin/bash

# DMXtoBeamer command-line arguments.
# Edit this file or create /mnt/config/dmxtobeamer.conf.
#
# Example only: uncomment and adapt to your application's
# actual command-line options.
#
# DMXTOBEAMER_ARGS=("--config" "/mnt/config/config.json")

DMXTOBEAMER_ARGS=()
EOF

# ============================================================
# Application startup wrapper
# ============================================================

cat > "$BIN_DIR/dmxtobeamer-start" <<'EOF'
#!/bin/bash
set -Eeuo pipefail

APP_DIR="/opt/dmxtobeamer"

# Default YAML configuration.
CONFIG_FILE="/etc/dmxtobeamer/config.yaml"

# Prefer the USB YAML configuration.
if [[ -f /mnt/config/config.yaml ]]; then
    CONFIG_FILE="/mnt/config/config.yaml"
fi

# Default launch argument configuration.
ARGS_FILE="/etc/dmxtobeamer/launch.conf"

# Prefer USB launch arguments when available.
if [[ -f /mnt/config/launch.conf ]]; then
    ARGS_FILE="/mnt/config/launch.conf"
fi

# Initialize launch arguments.
DMXTOBEAMER_ARGS=()

# Load launch arguments.
if [[ -f "$ARGS_FILE" ]]; then
    source "$ARGS_FILE"
fi

# Validate launch argument array.
if ! declare -p DMXTOBEAMER_ARGS 2>/dev/null | grep -q '^declare -a'; then
    echo "ERROR: DMXTOBEAMER_ARGS must be a Bash array."
    exit 1
fi

# Pass the selected YAML configuration to the Go application.
export DMXTOBEAMER_CONFIG="$CONFIG_FILE"

echo "Starting DMXtoBeamer"
echo "YAML configuration: $CONFIG_FILE"
echo "Launch arguments: $ARGS_FILE"
echo "Argument count: ${#DMXTOBEAMER_ARGS[@]}"

# Force the active virtual terminal to TTY1
/usr/bin/chvt 1
sleep 1

exec /usr/bin/cage -s -- "$APP_DIR/client" "${DMXTOBEAMER_ARGS[@]}"
EOF

chmod 0755 "$BIN_DIR/dmxtobeamer-start"

# ============================================================
# Kiosk user and system configuration hook
# ============================================================

mkdir -p config/hooks/live

cat > config/hooks/live/0100-dmxtobeamer.hook.chroot <<'EOF'
#!/bin/bash
set -e

# Create the kiosk user with a predictable UID and GID.
if ! getent group kiosk >/dev/null 2>&1; then
    groupadd --gid 1000 kiosk
fi
# Ensure all groups required by the kiosk service exist.
for group in video input render audio dialout seat; do
    if ! getent group "$group" >/dev/null 2>&1; then
        groupadd --system "$group"
    fi
done
if ! id kiosk >/dev/null 2>&1; then
    useradd \
        --uid 1000 \
        --gid 1000 \
        --create-home \
        --shell /bin/bash \
        kiosk
fi

# Grant hardware, journal, and seat access.
for group in video input render audio dialout plugdev systemd-journal sudo seat; do
    if getent group "$group" >/dev/null 2>&1; then
        usermod -aG "$group" kiosk
    fi
done

# Configure passwordless sudo.
mkdir -p /etc/sudoers.d

echo 'kiosk ALL=(ALL:ALL) NOPASSWD: ALL' > /etc/sudoers.d/kiosk

chmod 0440 /etc/sudoers.d/kiosk

visudo -cf /etc/sudoers.d/kiosk

# Configure German locale.
if ! grep -q '^de_DE.UTF-8 UTF-8' /etc/locale.gen; then
    echo 'de_DE.UTF-8 UTF-8' >> /etc/locale.gen
fi

locale-gen
update-locale LANG=de_DE.UTF-8

# Configure German keyboard.
cat > /etc/default/keyboard <<'KEYBOARD'
XKBMODEL="pc105"
XKBLAYOUT="de"
XKBVARIANT=""
XKBOPTIONS=""
BACKSPACE="guess"
KEYBOARD

# Create CONFIG mount point.
mkdir -p /mnt/config

# Enable required services.
systemctl enable NetworkManager.service
systemctl enable seatd.service
systemctl enable dmxtobeamer-config.service
systemctl enable dmxtobeamer.service
systemctl enable dmxtobeamer-logs.service
systemctl enable dmxtobeamer-system-logs.service

# Enable configuration terminal.
systemctl enable getty@tty3.service

# Disable getty services that conflict with our terminals.
systemctl disable getty@tty1.service || true
systemctl disable getty@tty2.service || true
systemctl disable getty@tty9.service || true
EOF

chmod 0755 config/hooks/live/0100-dmxtobeamer.hook.chroot

# ============================================================
# CONFIG partition mount service
# ============================================================

cat > "$SYSTEMD_DIR/dmxtobeamer-config.service" <<'EOF'
[Unit]
Description=Mount DMXtoBeamer CONFIG partition
After=local-fs.target
Before=dmxtobeamer.service

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=/bin/bash -c 'mkdir -p /mnt/config; for i in $(seq 1 30); do if [ -e /dev/disk/by-label/CONFIG ]; then mount -t vfat -o rw,uid=1000,gid=1000,umask=0022 /dev/disk/by-label/CONFIG /mnt/config && exit 0; fi; sleep 1; done; exit 0'
ExecStop=/bin/umount /mnt/config

[Install]
WantedBy=multi-user.target
EOF

# ============================================================
# DMXtoBeamer logs on TTY2
# ============================================================

cat > "$SYSTEMD_DIR/dmxtobeamer-logs.service" <<'EOF'
[Unit]
Description=DMXtoBeamer Application Logs on TTY2
After=systemd-journald.service dmxtobeamer.service
Conflicts=getty@tty2.service

[Service]
Type=simple
ExecStart=/usr/bin/journalctl -b -f -n 100 -u dmxtobeamer.service -o short-iso --no-pager

TTYPath=/dev/tty2
StandardInput=tty-force
StandardOutput=tty
StandardError=journal

TTYReset=yes
TTYVHangup=no
TTYVTDisallocate=no

Restart=always
RestartSec=2

[Install]
WantedBy=multi-user.target
EOF

# ============================================================
# Full system journal on TTY9
# ============================================================

cat > "$SYSTEMD_DIR/dmxtobeamer-system-logs.service" <<'EOF'
[Unit]
Description=Full System Journal on TTY9
After=systemd-journald.service
Conflicts=getty@tty9.service

[Service]
Type=simple
ExecStart=/usr/bin/journalctl -b -f -n 100 -o short-iso --no-pager

TTYPath=/dev/tty9
StandardInput=tty-force
StandardOutput=tty
StandardError=journal

TTYReset=yes
TTYVHangup=no
TTYVTDisallocate=no

Restart=always
RestartSec=2

[Install]
WantedBy=multi-user.target
EOF

# ============================================================
# DMXtoBeamer kiosk on TTY1
# ============================================================

cat > "$SYSTEMD_DIR/dmxtobeamer.service" <<'EOF'
[Unit]
Description=DMXtoBeamer Kiosk
After=systemd-user-sessions.service systemd-logind.service NetworkManager.service seatd.service dmxtobeamer-config.service
Wants=NetworkManager.service seatd.service
Conflicts=getty@tty1.service

[Service]
Type=simple

User=kiosk
Group=kiosk
SupplementaryGroups=video input render audio dialout seat

WorkingDirectory=/opt/dmxtobeamer

Environment=HOME=/home/kiosk
Environment=USER=kiosk
Environment=LANG=de_DE.UTF-8
Environment=XDG_RUNTIME_DIR=/run/dmxtobeamer
Environment=LIBSEAT_BACKEND=seatd

RuntimeDirectory=dmxtobeamer
RuntimeDirectoryMode=0700

TTYPath=/dev/tty1
StandardInput=tty-force
StandardOutput=journal
StandardError=journal

TTYReset=yes
TTYVHangup=no
TTYVTDisallocate=no

ExecStart=/bin/bash /usr/local/bin/dmxtobeamer-start

Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
EOF

# ============================================================
# Configuration terminal on TTY3
# ============================================================

mkdir -p config/includes.chroot/etc/systemd/system/getty@tty3.service.d

cat > config/includes.chroot/etc/systemd/system/getty@tty3.service.d/autologin.conf <<'EOF'
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin kiosk --noclear --keep-baud %I 115200,38400,9600 linux
TTYVHangup=no
TTYVTDisallocate=no
EOF

# ============================================================
# F3 configuration helper
# ============================================================

cat > "$BIN_DIR/dmx-config" <<'EOF'
#!/bin/bash
set -euo pipefail

CONFIG_DIR="/mnt/config"

# Fall back to local ISO configuration if USB is unavailable.
if ! mountpoint -q "$CONFIG_DIR"; then
    echo "USB CONFIG partition is not mounted."
    CONFIG_DIR="/etc/dmxtobeamer"
fi

while true; do
    clear

    echo "======================================"
    echo "       DMXtoBeamer Configuration"
    echo "======================================"
    echo
    echo "Configuration directory: $CONFIG_DIR"
    echo
    echo "1) Edit YAML configuration"
    echo "2) Edit launch arguments"
    echo "3) Show current configuration"
    echo "4) Restart DMXtoBeamer Engine"
    echo "5) Exit"
    echo

    read -rp "Select option: " choice

    case "$choice" in
        1)
            sudo nano "$CONFIG_DIR/config.yaml"
            ;;

        2)
            sudo nano "$CONFIG_DIR/launch.conf"
            ;;

        3)
            echo
            echo "----- config.yaml -----"
            cat "$CONFIG_DIR/config.yaml" 2>/dev/null || echo "Not found"

            echo
            echo "----- launch.conf -----"
            cat "$CONFIG_DIR/launch.conf" 2>/dev/null || echo "Not found"

            echo
            read -rp "Press Enter to continue..."
            ;;

        4)
            echo
            echo "Restarting DMXtoBeamer Engine..."
            sudo systemctl restart dmxtobeamer.service
            echo "Restart command executed."
            sleep 2
            ;;

        5)
            exit 0
            ;;

        *)
            echo "Invalid selection."
            sleep 1
            ;;
    esac
done
EOF

chmod 0755 "$BIN_DIR/dmx-config"

# ============================================================
# Build ISO
# ============================================================

echo
echo "=========================================="
echo " Building Debian 13 Live ISO"
echo "=========================================="

cd "$WORK_DIR" || exit 1

# Clean Live Build artifacts and rebuild ISO to ensure consistency
lb bootstrap
lb build || {
    echo "ERROR: Live Build failed. See logs above."
    exit 1
}

# ============================================================
# Copy ISO
# ============================================================

BUILT_ISO="$WORK_DIR/live-image-amd64.hybrid.iso"

if [[ ! -f "$BUILT_ISO" ]]; then
    echo "ERROR: ISO was not generated:"
    echo "$BUILT_ISO"
    exit 1
fi

cp "$BUILT_ISO" "$ISO_OUTPUT"

# ============================================================
# Finished
# ============================================================

echo
echo "=========================================="
echo " BUILD COMPLETE"
echo "=========================================="

echo "ISO:  $ISO_OUTPUT"
echo "Size: $(du -h "$ISO_OUTPUT" | cut -f1)"
