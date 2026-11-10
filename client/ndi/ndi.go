// Package ndi provides NDI® video output for the client without linking or
// bundling the NDI SDK. The NDI runtime (libndi.so) is loaded dynamically at
// runtime via dlopen, which keeps the binary open source and usable on
// machines without NDI installed.
//
// The handful of C structs and entry points bound here are re-implemented
// against the documented public NDI C API (see https://docs.ndi.video);
// nothing from the NDI SDK itself is vendored or distributed with this
// project. Users must install the NDI runtime themselves, available from
// https://ndi.video/tools/.
//
// NDI® is a registered trademark of Vizrt NDI AB.
//
//go:build ndi
package ndi

import (
	"fmt"
	"log"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// init prints the NDI® startup notice required by the NDI SDK license terms
// (https://docs.ndi.video/all/developing-with-ndi/sdk/licensing). Because this
// file is gated behind the `ndi` build tag, the notice is shown exactly on
// builds that include NDI support, once at program start — regardless of
// whether --ndi/--ndi-only is used.
func init() {
	log.Println("This build includes NDI® support. NDI® is a registered trademark of Vizrt NDI AB — https://ndi.video/. " +
		"The NDI runtime is not bundled with this application and must be installed separately (https://ndi.video/tools/).")
}

// FourCC values (little-endian four-character codes from Processing.NDI.Structs.h).
const (
	FourCCTypeUYVY = 0x59565955 // 'UYVY' — 4:2:2 packed, 16 bits per pixel
	FourCCTypeBGRA = 0x41524742 // 'BGRA' — uncompressed 8-bit BGRA
	FourCCTypeBGRX = 0x58524742 // 'BGRX' — uncompressed 8-bit BGRX
)

// Frame format types (NDIlib_frame_format_type_e).
const (
	FrameFormatTypeProgressive = 1
)

// libraryNames are the candidate shared-library names/paths tried in order,
// matching the NDI 6/5/4 Linux runtime installs. The runtime is not
// distributed with this project and is never installed by us — see
// https://ndi.video/tools/.
var libraryNames = []string{
	"libndi.so.6",
	"libndi.so.5",
	"libndi.so.4",
	"libndi.so",
	"Processing.NDI.Lib.x64.so",
}

// entryPoints holds the NDI entry points we need.
type entryPoints struct {
	initialize func() bool
	destroy    func()
	sendCreate func(pCreate *sendCreateT) unsafe.Pointer
	sendDestroy func(instance unsafe.Pointer)
	sendVideoV2 func(instance unsafe.Pointer, frame *videoFrameV2T)
	sendNoConnections func(instance unsafe.Pointer, timeoutMs uint32) int32
}

// videoFrameV2T mirrors NDIlib_video_frame_v2_t (Processing.NDI.Structs.h).
type videoFrameV2T struct {
	xres               int32
	yres               int32
	FourCC             int32
	frameRateN         int32
	frameRateD         int32
	pictureAspectRatio float32
	frameFormatType    int32
	timecode           int64
	data               unsafe.Pointer
	lineStrideInBytes  int32
	metadata           *byte // const char* — NULL = none
	reserved           uint32
	timestamp          int64
}

// sendCreateT mirrors NDIlib_send_create_t (Processing.NDI.Structs.h).
type sendCreateT struct {
	ndiName    *byte // const char*
	groups     *byte // const char*
	clockVideo bool
	clockAudio bool
}

// Library holds the loaded NDI runtime.
type Library struct {
	sym entryPoints
}

var (
	libOnce sync.Once
	loaded  *Library
	loadErr error
)

// Load locates and initializes the NDI runtime. It is safe to call
// multiple times; the result is cached. Returns nil (with an error) when
// the runtime is not installed.
//
// The NDI runtime is not distributed with this project — it must be
// installed by the user, available from https://ndi.video/tools/.
func Load() (*Library, error) {
	libOnce.Do(func() {
		lib, err := loadLibrary()
		if err != nil {
			loadErr = err
			return
		}
		if !lib.sym.initialize() {
			loadErr = fmt.Errorf("NDIlib_initialize failed (NDI runtime found, but initialization returned false)")
			return
		}
		loaded = lib
	})
	return loaded, loadErr
}

// Shutdown destroys the NDI library context. Safe to call even if Load
// failed or was never called.
func Shutdown() {
	if loaded != nil {
		loaded.sym.destroy()
		loaded = nil
		loadErr = nil
		libOnce = sync.Once{}
	}
}

// loadLibrary tries each candidate name and resolves the entry points we
// need directly by symbol name through dlsym — stable across NDI 4/5/6.
func loadLibrary() (*Library, error) {
	var lastErr error
	for _, name := range libraryNames {
		h, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_LOCAL)
		if err != nil {
			lastErr = err
			continue
		}
		lib, err := bindLibrary(h)
		if err == nil {
			return lib, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("NDI runtime library not found (tried %v): %w; "+
		"please install the NDI runtime from https://ndi.video/tools/ "+
		"(NDI® is a registered trademark of Vizrt NDI AB)", libraryNames, lastErr)
}

// bindLibrary resolves the entry points we need from a loaded handle.
func bindLibrary(h uintptr) (*Library, error) {
	lib := &Library{}
	purego.RegisterLibFunc(&lib.sym.initialize, h, "NDIlib_initialize")
	purego.RegisterLibFunc(&lib.sym.destroy, h, "NDIlib_destroy")
	purego.RegisterLibFunc(&lib.sym.sendCreate, h, "NDIlib_send_create")
	purego.RegisterLibFunc(&lib.sym.sendDestroy, h, "NDIlib_send_destroy")
	purego.RegisterLibFunc(&lib.sym.sendVideoV2, h, "NDIlib_send_send_video_v2")
	purego.RegisterLibFunc(&lib.sym.sendNoConnections, h, "NDIlib_send_get_no_connections")
	return lib, nil
}

// Sender is a live NDI sending instance.
type Sender struct {
	instance unsafe.Pointer
	sym      entryPoints
	name     string
}

// NewSender creates an NDI sender. groups may be "".
func (l *Library) NewSender(name string, groups string) (*Sender, error) {
	if l == nil {
		return nil, fmt.Errorf("NDI not loaded")
	}
	sc := sendCreateT{
		ndiName:    cstr(name),
		clockVideo: true,
	}
	if groups != "" {
		sc.groups = cstr(groups)
	}
	inst := l.sym.sendCreate(&sc)
	if inst == nil {
		return nil, fmt.Errorf("NDIlib_send_create failed for %q", name)
	}
	return &Sender{instance: inst, sym: l.sym, name: name}, nil
}

// Name returns the NDI source name.
func (s *Sender) Name() string {
	if s == nil {
		return ""
	}
	return s.name
}

// SendVideo submits a frame. With NDIlib_send_send_video_v2 the call
// returns once the library has copied the frame data, so the buffer may be
// reused for the next frame after it returns.
func (s *Sender) SendVideo(xres, yres, fourCC, lineStride, frameRateN, frameRateD int32, data unsafe.Pointer) {
	if s == nil || s.instance == nil {
		return
	}
	frame := videoFrameV2T{
		xres:              xres,
		yres:              yres,
		FourCC:            fourCC,
		frameRateN:        frameRateN,
		frameRateD:        frameRateD,
		frameFormatType:   FrameFormatTypeProgressive,
		timecode:          0, // 0 = auto-generate timecode
		data:              data,
		lineStrideInBytes: lineStride,
	}
	s.sym.sendVideoV2(s.instance, &frame)
}

// ConnectionCount returns how many receivers are currently connected.
func (s *Sender) ConnectionCount() int {
	if s == nil || s.instance == nil {
		return 0
	}
	return int(s.sym.sendNoConnections(s.instance, 0))
}

// Close destroys the sender instance.
func (s *Sender) Close() {
	if s == nil || s.instance == nil {
		return
	}
	s.sym.sendDestroy(s.instance)
	s.instance = nil
}

// cstr allocates a NUL-terminated copy of s. The backing array is kept
// alive by the returned pointer.
func cstr(s string) *byte {
	b := append([]byte(s), 0)
	return &b[0]
}
