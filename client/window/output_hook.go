// NDI® is a registered trademark of Vizrt NDI AB, https://ndi.video/
//
//go:build ndi

package window

import (
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/config"
	"technikflg.com/dmxToProjector/ndi"
)

// This file is compiled only in builds that include NDI support (build tag
// `ndi`). Builds without the tag use output_hook_noop.go instead.

// outputManager is the process-wide streaming output manager. It is created
// by InitOutput when --ndi or --ndi-only is active and receives one frame per
// rendered frame after EndDrawing.
var outputManager *ndi.Manager

// outputOnlyMode records headless mode (no visible window) for the render loop.
var outputOnlyMode bool

// outputTargetW/outputTargetH is the requested streaming output resolution
// (0 = send the window render resolution unchanged). Set by InitOutput from
// --ndi-width and --ndi-height, which now apply in windowed mode too.
var (
	outputTargetW int32
	outputTargetH int32
)

// outputFirstFrameLogged makes the one-time "first frame sent" hint log fire once.
var outputFirstFrameLogged bool

// InitOutput prepares the streaming output. Call before InitWindowEx. When
// onlyMode is true, no visible window will be opened later (headless
// --ndi-only mode). outW/outH is the output resolution (0 = window render
// resolution).
func InitOutput(name, groups string, onlyMode bool, outW, outH int32) {
	outputOnlyMode = onlyMode
	outputTargetW = outW
	outputTargetH = outH
	outputManager = ndi.NewManager(name, groups)
}

// OutputManager returns the manager (nil when the output is not active).
func OutputManager() *ndi.Manager {
	return outputManager
}

// enableOutput creates the streaming sender once the render resolution is
// known. Called from the render path on the GL thread.
func enableOutput(w, h int32) {
	if outputManager == nil || outputManager.Enabled() {
		return
	}
	// config.Cfg may be nil in production (configuration arrives via
	// WebSocket); fall back to the same 60 fps default as window init.
	fps := int32(60)
	if cfg := config.Cfg; cfg != nil && cfg.Window.MaxFPS > 0 {
		fps = cfg.Window.MaxFPS
	}
	if err := outputManager.Enable(w, h, fps, outputTargetW, outputTargetH); err != nil {
		// NDI runtime missing: windowed mode keeps running without the
		// stream, headless mode cannot fulfill its purpose and exits in
		// main().
		log.Printf("NDI® output disabled: %v", err)
		outputManager = nil
	}
}

// pushOutputFrame reads the finished framebuffer and hands it to the streaming
// sender. Must be called from the raylib/GL thread, right after
// rl.EndDrawing().
func pushOutputFrame() {
	if outputManager == nil {
		return
	}
	outputManager.PushFrameBGRA()
	if !outputFirstFrameLogged {
		outputFirstFrameLogged = true
		log.Println("NDI®: first frame sent — the source should now appear in NDI receivers (discovery may take a few seconds). " +
			"If no receiver connects: mDNS discovery uses UDP port 5353 (allow it in the firewall) or connect to the sender by IP address.")
	}
}

// ShutdownOutput closes the streaming sender and tears down the library.
func ShutdownOutput() {
	if outputManager != nil {
		outputManager.Shutdown()
		outputManager = nil
	}
}

// InitWindowHidden prepares rendering without showing a window, for headless
// --ndi-only mode. The window is created hidden at a fixed output resolution
// so the GPU rendering pipeline (and hardware encoding via the compositor)
// stays fully active.
func InitWindowHidden(WindowName string, width, height int32) {
	// A hidden GL window still needs a display connection on X11, so we cannot
	// run on a truly display-less machine. Give a clean error up front.
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		log.Println("--ndi-only: No display available (DISPLAY/WAYLAND_DISPLAY unset); a hidden GL window is still needed to render with the GPU and will likely fail.")
		log.Println("--ndi-only: Run the client on a machine with a running X11/Wayland session.")
	}

	WindowWidth = int(width)
	WindowHeight = int(height)

	// The production client does not load config.yaml (configuration comes
	// from the server via WebSocket), so Cfg may be nil here. Fall back to
	// the defaults used elsewhere.
	maxFPS := int32(60)
	vsync := false
	if cfg := config.Cfg; cfg != nil {
		if cfg.Window.MaxFPS > 0 {
			maxFPS = cfg.Window.MaxFPS
		}
		vsync = cfg.Window.Vsync
	}

	rl.SetConfigFlags(rl.FlagWindowHidden)
	rl.InitWindow(width, height, WindowName+" Animation Client (NDI)")
	if !rl.IsWindowReady() {
		log.Println("--ndi-only: Failed to initialize the rendering window — cannot run headless.")
		return
	}

	initWindowCommon(maxFPS, vsync, width, height)
}

// isWindowHiddenMode reports whether headless streaming mode is active.
func isWindowHiddenMode() bool {
	return outputOnlyMode
}
