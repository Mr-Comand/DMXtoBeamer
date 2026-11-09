// No-op NDI stubs for builds compiled without the `ndi` build tag (for
// example the Windows client, where the NDI runtime cannot be loaded with
// purego — purego.Dlopen is unavailable on Windows).
//
// With the --ndi / --ndi-only flags on such a build, the client logs a
// hint and keeps running in normal windowed mode; headless --ndi-only mode
// exits, since it cannot fulfill its purpose without NDI.
//
// NDI® is a registered trademark of Vizrt NDI AB, https://ndi.video/
//
//go:build !ndi

package window

import "log"

// ndiOnlyMode records headless mode (no visible window) for the render loop.
var ndiOnlyMode bool

// InitNDI prepares NDI output in NDI-enabled builds. In builds without the
// `ndi` tag it only records onlyMode: --ndi-only can never work here, so
// IsNDIOnly reports it and the caller exits headless mode early.
func InitNDI(name, groups string, onlyMode bool) {
	ndiOnlyMode = onlyMode
	if onlyMode {
		log.Println("This build was compiled without NDI® support — headless --ndi-only mode is not available. Rebuild with the 'ndi' build tag (see LINUX.md).")
	} else {
		log.Printf("This build was compiled without NDI® support — ignoring NDI® options (--ndi --ndi-name=%q --ndi-groups=%q). Rebuild with the 'ndi' build tag (see LINUX.md).", name, groups)
	}
}

// NDIManager returns nil so callers behave as if NDI failed to start.
func NDIManager() *struct{} {
	return nil
}

// IsNDIOnly reports whether headless --ndi-only mode is active.
func IsNDIOnly() bool {
	return ndiOnlyMode
}

// enableNDI is a no-op in builds without the `ndi` tag.
func enableNDI(w, h int32) {}

// pushNDIFrame is a no-op in builds without the `ndi` tag.
func pushNDIFrame() {}

// ShutdownNDI is a no-op in builds without the `ndi` tag.
func ShutdownNDI() {}
