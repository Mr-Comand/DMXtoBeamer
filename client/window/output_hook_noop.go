//go:build !ndi

package window

// Silent no-op output stubs for builds compiled without the `ndi` build tag
// (for example the Windows client, where the NDI runtime cannot be loaded
// with purego — purego.Dlopen is unavailable on Windows). This file contains
// no references to the optional streaming backend, so non-streaming builds
// carry no trace of it.

// InitOutput is a no-op in builds without the `ndi` tag.
func InitOutput(name, groups string, onlyMode bool, outW, outH int32) {}

// OutputManager returns nil so callers behave as if the output failed to
// start.
func OutputManager() *struct{} { return nil }

// enableOutput is a no-op in builds without the `ndi` tag.
func enableOutput(w, h int32) {}

// pushOutputFrame is a no-op in builds without the `ndi` tag.
func pushOutputFrame() {}

// ShutdownOutput is a no-op in builds without the `ndi` tag.
func ShutdownOutput() {}

// InitWindowHidden is a no-op in builds without the `ndi` tag.
func InitWindowHidden(WindowName string, width, height int32) {}

// isWindowHiddenMode is always false in builds without the `ndi` tag.
func isWindowHiddenMode() bool { return false }
