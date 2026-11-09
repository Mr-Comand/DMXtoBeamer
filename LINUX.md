# Linux / Raspberry Pi Setup

This project can run natively on Linux, including on a Raspberry Pi.
The Windows build instructions in `client/build.sh` are unchanged.

## Requirements

### Server (`server/`)
- Go 1.21+
- No extra libraries (pure Go HTTP/WebSocket server)

```sh
cd server
go build -o dmxtobeamerserver .
./dmxtobeamerserver
```

The server listens on `http://<host>:8080` by default.

### Client (`client/`)
- Go 1.21+ (cgo enabled)
- PortAudio + ALSA development headers
- A display (raylib window; works on X11 and Wayland)

Debian / Raspberry Pi OS (Bookworm or later):

```sh
sudo apt install golang portaudio19-dev libasound2-dev \
  libx11-dev libxkbcommon-dev libxcursor-dev libxrandr-dev \
  libxinerama-dev libxi-dev libgl-dev libglu1-mesa-dev
```

Fedora / other RPM distributions:

```sh
sudo dnf install golang portaudio-devel alsa-lib-devel \
  libX11-devel libxkbcommon-devel libXcursor-devel libXrandr-devel \
  libXinerama-devel libXi-devel mesa-libGL-devel \
  mesa-libGLU-devel wayland-devel
```

For audio capture the system needs a working ALSA input device.
On desktop systems that use PipeWire or PulseAudio, install
`pipewire-alsa` / `pulseaudio-alsa` (usually present by default).

## Building the client

```sh
cd client
./build-linux.sh          # native build for this machine
./build-linux.sh arm64    # cross-compile for Raspberry Pi OS 64-bit
./build-linux.sh arm      # cross-compile for Raspberry Pi OS 32-bit (armv7)
```

Cross-compiling needs the matching toolchain:

- Debian/Ubuntu: `sudo apt install gcc-aarch64-linux-gnu` (64-bit)
  or `gcc-arm-linux-gnueabihf` (32-bit)
- Fedora: `sudo dnf install gcc-aarch64-linux-gnu` or `gcc-arm-linux-gnu`

The result is placed in `client/build/` together with the `shaders/`
directory, which must be kept next to the binary (it is loaded at
runtime from the working directory).

The client is built with the raylib `x11` build tag, which compiles
GLFW with X11 support only. This avoids needing Wayland development
headers at build time; on a Wayland desktop the binary still runs
through XWayland.

## Running the client

```sh
cd client/build
./client ws://<server-host>:8080 <clientName>
```

On startup the client creates a `config.yaml` in its working directory.

### Audio input notes (Linux)

Linux audio devices behave differently from Windows/WASAPI:

- The PortAudio default device maps to the ALSA `default` PCM, which
  fails to open on some setups (notably PipeWire). The client detects
  this and falls back to the `sysdefault` ALSA device automatically.
- On Linux the client always opens **all input channels** of the chosen
  device and downmixes them to mono. Many ALSA hardware devices refuse
  to open with a single channel, so this avoids `Device unavailable`
  errors.
- If the fallback still fails, the device can be selected by name in
  `config.yaml` (`audioDevice`), e.g. `sysdefault` or `hw:0,0`.

## NDI® output

NDI® is a registered trademark of Vizrt NDI AB (https://ndi.video/).

The client can additionally (or exclusively) stream its rendered output
to the network as an NDI® source. The NDI runtime is **not** bundled:
install it separately from https://ndi.video/tools/ (the runtime library
is found automatically as `libndi.so.6`/`.5`/`.4` in the system loader
path, e.g. `/usr/local/lib` — run `sudo ldconfig` after installing).

NDI options:

```sh
# Stream to NDI in addition to the normal window:
./client --ndi ws://<server-host>:8080 <clientName>

# Headless NDI-only mode (no visible window, still GPU-rendered):
./client --ndi-only ws://<server-host>:8080 <clientName>

# Optional parameters:
./client --ndi-name="Stage Left" --ndi-groups=GroupA,GroupB \
         --ndi-width=1920 --ndi-height=1080 --ndi-only ws://host:8080 Name
```

- `--ndi`: stream the rendered output to NDI in addition to the window.
- `--ndi-only`: stream to NDI only — no visible window is shown.
  Rendering still uses the GPU (hardware accelerated OpenGL). The
  hidden render window requires a running X11/Wayland session, so this
  mode cannot run on a display-less server; the render resolution can
  be chosen with `--ndi-width`/`--ndi-height` (default 1280x720).
- `--ndi-name`: NDI source name (default: `DMXToBeamer-<clientName>`).
- `--ndi-groups`: comma separated NDI discovery groups.
- `--ndi-width` / `--ndi-height`: render resolution in `--ndi-only`
  mode.

Frames are sent uncompressed (UYVY 4:2:2); NDI handles compression and
transport itself, so no third-party video codec licensing applies.
Recipients discover the source with any NDI receiver, e.g. NDI Studio
Monitor (https://ndi.video/tools/).

## Raspberry Pi specifics

- Raspberry Pi OS 64-bit (Bookworm) uses PipeWire by default; the
  client works with it through the ALSA plugin layer.
- Run the client from a desktop session (it opens a raylib window).
- For headless test servers, only `server/` is needed — it has no
  graphical dependencies.
