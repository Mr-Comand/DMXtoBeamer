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

// InitNDI prepares NDI output. Call before InitWindowEx. When onlyMode is
// true, no visible window will be opened later (headless --ndi-only mode).
func InitNDI(name, groups string, onlyMode bool) {
	ndiOnlyMode = onlyMode
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
	if err := ndiManager.Enable(w, h, fps); err != nil {
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
}

// ShutdownNDI closes the NDI sender and tears down the NDI library.
func ShutdownNDI() {
	if ndiManager != nil {
		ndiManager.Shutdown()
		ndiManager = nil
	}
}
