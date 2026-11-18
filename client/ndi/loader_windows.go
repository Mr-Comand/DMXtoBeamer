//go:build ndi && windows

package ndi

import "golang.org/x/sys/windows"

var libraryNames = []string{
	"Processing.NDI.Lib.x64.dll",
}

func openLibrary(name string) (uintptr, error) {
	handle, err := windows.LoadLibrary(name)
	return uintptr(handle), err
}
