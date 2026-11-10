// NDI® is a registered trademark of Vizrt NDI AB, https://ndi.video/
//
//go:build ndi

package window

import (
	"log"

	"technikflg.com/dmxToProjector/config"
	"technikflg.com/dmxToProjector/ndi"
)

// This file is compiled only in builds that include NDI support (build tag
// `ndi`). Builds without the tag use ndi_hook_noop.go instead.

// ndiManager is the process-wide NDI output manager. It is created by
// InitNDI when --ndi or --ndi-only is active and receives one frame per
// rendered frame after EndDrawing.
var ndiManager *ndi.Manager

// ndiOnlyMode records headless mode (no visible window) for the render loop.
var ndiOnlyMode bool

// ndiTargetW/ndiTargetH is the requested NDI output resolution (0 = send the
// window render resolution unchanged). Set by InitNDI from --ndi-width and
// --ndi-height, which now apply in windowed mode too.
var (
	ndiTargetW int32
	ndiTargetH int32
)

// ndiFirstFrameLogged makes the one-time "first frame sent" hint log fire once.
var ndiFirstFrameLogged bool

// InitNDI prepares NDI output. Call before InitWindowEx. When onlyMode is
// true, no visible window will be opened later (headless --ndi-only mode).
// outW/outH is the NDI output resolution (0 = window render resolution).
func InitNDI(name, groups string, onlyMode bool, outW, outH int32) {
	ndiOnlyMode = onlyMode
	ndiTargetW = outW
	ndiTargetH = outH
	ndiManager = ndi.NewManager(name, groups)
}

// NDIManager returns the manager (nil when NDI is not active).
func NDIManager() *ndi.Manager {
	return ndiManager
}

// IsNDIOnly reports whether headless --ndi-only mode is active.
func IsNDIOnly() bool {
	return ndiOnlyMode
}

// enableNDI creates the NDI sender once the render resolution is known.
// Called from the render path on the GL thread.
func enableNDI(w, h int32) {
	if ndiManager == nil || ndiManager.Enabled() {
		return
	}
	// config.Cfg may be nil in production (configuration arrives via
	// WebSocket); fall back to the same 60 fps default as window init.
	fps := int32(60)
	if cfg := config.Cfg; cfg != nil && cfg.Window.MaxFPS > 0 {
		fps = cfg.Window.MaxFPS
	}
	if err := ndiManager.Enable(w, h, fps, ndiTargetW, ndiTargetH); err != nil {
		// NDI runtime missing: windowed mode keeps running without NDI,
		// headless mode cannot fulfill its purpose and exits in main().
		log.Printf("NDI® output disabled: %v", err)
		ndiManager = nil
	}
}

// pushNDIFrame reads the finished framebuffer and hands it to the NDI sender.
// Must be called from the raylib/GL thread, right after rl.EndDrawing().
func pushNDIFrame() {
	if ndiManager == nil {
		return
	}
	ndiManager.PushFrameBGRA()
	if !ndiFirstFrameLogged {
		ndiFirstFrameLogged = true
		log.Println("NDI®: first frame sent — the source should now appear in NDI receivers (discovery may take a few seconds). " +
			"If no receiver connects: mDNS discovery uses UDP port 5353 (allow it in the firewall) or connect to the sender by IP address.")
	}
}

// ShutdownNDI closes the NDI sender and tears down the NDI library.
func ShutdownNDI() {
	if ndiManager != nil {
		ndiManager.Shutdown()
		ndiManager = nil
	}
}
