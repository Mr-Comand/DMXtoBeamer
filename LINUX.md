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
sudo apt install golang portaudio19-dev libasound2-dev
```

Fedora / other RPM distributions:

```sh
sudo dnf install golang portaudio-devel alsa-lib-devel
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

## Raspberry Pi specifics

- Raspberry Pi OS 64-bit (Bookworm) uses PipeWire by default; the
  client works with it through the ALSA plugin layer.
- Run the client from a desktop session (it opens a raylib window).
- For headless test servers, only `server/` is needed — it has no
  graphical dependencies.
