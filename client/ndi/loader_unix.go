//go:build ndi && !windows

package ndi

import "github.com/ebitengine/purego"

var libraryNames = []string{
	"libndi.so.6",
	"libndi.so.5",
	"libndi.so.4",
	"libndi.so",
	"Processing.NDI.Lib.x64.so",
}

func openLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_LOCAL)
}
