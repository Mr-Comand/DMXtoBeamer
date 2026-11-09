# NDI® Streaming — TODO

Add NDI® output to the Linux client: `--ndi` streams to NDI **in addition** to the
normal raylib window, `--ndi-only` streams NDI **instead of** opening a window
(headless). Both modes keep the normal raylib GPU rendering (hardware accelerated
OpenGL pipeline) — only the final framebuffer is additionally read back and sent.

> **NDI® is a registered trademark of Vizrt NDI AB.**
> <https://ndi.video/> — NDI SDK licensing terms:
> <https://docs.ndi.video/all/developing-with-ndi/sdk/licensing>

## Licensing compliance (docs.ndi.video/all/developing-with-ndi/sdk/licensing)

- [x] **Do not redistribute the NDI SDK / NDI libraries.** The client loads the
  NDI runtime (`libndi.so.6` / `libndi.so.5` / `libndi.so.4`) **dynamically at
  runtime** (`dlopen`) — nothing NDI-owned is bundled, linked or committed to
  this open-source repository.
- [x] **Do not install NDI libraries into the system path from our side.** The
  client only *searches* for an already installed runtime; it never installs or
  copies NDI files itself. Users install the NDI runtime from
  <https://ndi.video/tools/> themselves.
- [x] **Link to ndi.video** where NDI is used/selected: shown in the client
  startup log, in the headless `--help` output, and documented in README.md and
  LINUX.md next to the NDI feature description.
- [x] **Trademark attribution** "NDI® is a registered trademark of Vizrt NDI AB"
  is printed near the first use of the word NDI (startup log) and included in
  README.md, LINUX.md and this file.
- [x] **No NDI tools distribution.** Documentation points users to
  <https://ndi.video/tools/> to obtain the NDI runtime and NDI tools.
- [x] **H.264/H.265/AAC caution.** The client sends uncompressed 4:2:2
  video (RGBA→UYVY conversion done in the client). No H.264/H.265/AAC
  encoding is used, so no additional codec licensing applies.

## Design

- [x] Design decision: read the final GPU framebuffer with
  `rl.LoadImageFromScreen` (RGBA readback via `glReadPixels`) and convert to
  UYVY (BT.601, 4:2:2) in the client, then send as
  `NDIlib_FourCC_video_type_UYVY` with `NDIlib_send_send_video_v2`, keeping
  the NDI library free to do its own compression (SpeedHQ) and transport.
- [x] Keep all raylib rendering unchanged; NDI output is a pure add-on
  ("tap") after `rl.EndDrawing()` — animations, shaders and window behavior
  are untouched.
- [x] Dynamic loading (`dlopen`) so the binary stays fully open source
  (BSD-2 style NDI SDK license allows the headers' API definitions to be
  re-implemented; no SDK source or binary is vendored) and runs on machines
  without NDI installed (window mode keeps working, NDI is just disabled).

## Implementation

- [x] New `ndi` branch from the `linux` branch for development.
- [x] Add NDI C-API bindings (`client/ndi/ndi.go`) — pure Go structs +
  `unsafe`/`syscall`-based dynamic loading via `purego` (no cgo needed):
  initialize, sender create, `NDIlib_send_send_video_v2`, destroy.
- [x] UYVY (4:2:2) converter from RGBA framebuffer data
  (`client/ndi/convert.go`) — BT.601 full-range conversion, OpenGL bottom-up
  rows flipped to top-down, odd widths cropped to even (UYVY requirement);
  covered by unit tests in `client/ndi/convert_test.go`.
- [x] Frame source hook in `window.Render()` — after `rl.EndDrawing()` read
  the framebuffer with `glReadPixels` (BGRA) and feed it to the NDI sender
  when enabled (`client/window/ndi_hook.go`).
- [x] Flag parsing: support `--ndi` (stream NDI in addition to the window),
  `--ndi-only` (stream NDI only, no window = headless), `--ndi-name=NAME`
  (NDI source name, default `DMXToBeamer-<ClientId>`), `--ndi-groups=LIST`
  (NDI discovery groups) and `--ndi-width`/`--ndi-height` (headless render
  resolution). Unknown flags are warned about so positional URL/name parsing
  keeps working.
  In `--ndi-only` mode no visible raylib window is created; the renderer
  still runs with hardware acceleration via a hidden raylib window
  (`rl.SetConfigFlags(rl.FlagWindowHidden)`, fixed 1280×720 default
  resolution) — the full OpenGL pipeline stays active ("headless hardware
  rendering").
- [x] `--ndi-only` mode: hidden raylib window at a fixed internal resolution
  (1280×720 default, `--ndi-width`/`--ndi-height` to change), render loop
  unchanged (`window.InitWindowHeadless`), framebuffer read back and sent
  after every frame. Warns cleanly when no X11/Wayland display is available
  (a hidden GL window still needs a display connection).
- [x] Flag parsing for both `main.go` (dev build) and `prod.go`
  (production build) — shared helper `parseFlags()` in `cliflags.go`
  (build-tag compatible, no build tag) to avoid code duplication.
- [x] Graceful degradation: if `dlopen` of the NDI runtime fails, log a clear
  error, and — in `--ndi-only` mode — exit with an error (window mode just
  continues without NDI).
- [x] Clean shutdown: `ndi.Shutdown()` on program exit, sender destroyed
  before NDI init teardown.

## Build system

- [x] `client/build-linux.sh` builds the client binary with NDI support
  (purego/dlopen — no extra build dependencies, no NDI SDK needed at build
  time).
- [x] LINUX.md: NDI runtime installation notes (user installs NDI runtime
  from ndi.video/tools), licensing note, trademark attribution.
- [ ] README.md: NDI usage section (no root README.md exists yet — the Linux
  documentation is LINUX.md; skip unless a README is wanted).

## Documentation

- [x] Document flags and NDI source naming in LINUX.md.
- [x] Document that NDI runtime must be installed by the user (not bundled)
  and where to get it (ndi.video / ndi.video/tools), with the link requirement
  from the NDI SDK license.
- [x] Trademark footnote in LINUX.md (no root README.md exists — see above).

## Testing / verification

- [x] `go build ./...` and `go vet ./...` pass in `client/` (both tag sets:
  plain `x11` dev build and `production x11`).
- [x] Unit tests for the UYVY converter (`client/ndi/convert_test.go`):
  primary colors, row flip, chroma averaging, odd-width cropping — all pass.
- [x] `--help` prints the NDI option summary with the trademark attribution
  and ndi.video links.
- [x] Without NDI runtime installed: binary links fine (no NDI linkage at
  build time), window mode starts normally and logs a clear error explaining
  how to install the runtime; `--ndi-only` likewise logs and exits.
- [x] Clean SIGINT/SIGTERM shutdown in headless mode: signal handler sets a
  shutdown flag checked by the render loop; `end()` no longer blocks on a
  nil audio stop channel (`StopAudioProcessor` guards it); SIGTERM smoke
  test exits cleanly.
- [x] `InitWindowHeadless` hardened: nil-safe config fallback (production
  builds never load config.yaml), failed window init is detected and the
  render loop is skipped instead of crashing on raylib calls.
- [ ] Window mode (`--ndi`): window opens as before, NDI sender registers
  (visible in NDI Studio Monitor), frames flow — requires a machine with
  the NDI runtime installed.
- [ ] Headless mode (`--ndi-only`) on a machine with a GPU: NDI stream
  visible in NDI Studio Monitor, animations run — requires NDI runtime.
- [ ] Live end-to-end test with NDI Studio Monitor on real hardware
  (requires a machine with NDI runtime installed — left for manual
  verification by the user).
