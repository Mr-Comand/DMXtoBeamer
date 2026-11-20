//go:build ndi

package ndi

import (
	"fmt"
	"log"
	"sync"
	"time"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Manager owns the NDI sender lifecycle and receives finished frames from the
// render loop. It is created once and lives for the whole program.
type Manager struct {
	lib     *Library
	sender  *Sender
	name    string
	groups  string
	mu      sync.Mutex
	started bool

	// outW/outH is the requested NDI output resolution. 0 means "use the
	// window render resolution unchanged" (no scaling).
	outW int32
	outH int32

	// Frame scratch buffers (reused; SendVideo copies data before returning).
	bgraBuf  []byte
	uyvyBuf  []byte
	curW     int32
	curH     int32

	// fps bookkeeping for frame-rate metadata
	fpsN   int32
	fpsD   int32

	lastLog time.Time
	frames  uint64
}

// NewManager returns a manager without a live sender yet; Enable actually
// creates the NDI sender (after the render target resolution is known).
func NewManager(name, groups string) *Manager {
	return &Manager{name: name, groups: groups}
}

// Enable creates the NDI sender once the initial resolution is known.
// outW/outH is the NDI output resolution; 0 for either dimension means the
// frame is sent at the window render resolution without scaling.
// Returns an error when the NDI runtime is not available.
func (m *Manager) Enable(w, h int32, fps int32, outW, outH int32) error {
	if m == nil {
		return fmt.Errorf("ndi manager is nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.started {
		return nil
	}
	m.outW = even(outW)
	m.outH = outH
	lib, err := Load()
	if err != nil {
		return err
	}
	m.lib = lib
	m.fpsN, m.fpsD = fpsToFraction(fps)
	sender, err := lib.NewSender(m.name, m.groups)
	if err != nil {
		return err
	}
	m.sender = sender
	m.started = true
	if m.outW > 0 && m.outH > 0 {
		log.Printf("NDI® sender created: %q (window %dx%d, output %dx%d @ %d fps) — NDI® is a registered trademark of Vizrt NDI AB; get the NDI runtime and tools at https://ndi.video/tools/",
			sender.Name(), w, h, m.outW, m.outH, fps)
	} else {
		log.Printf("NDI® sender created: %q (%dx%d @ %d fps) — NDI® is a registered trademark of Vizrt NDI AB; get the NDI runtime and tools at https://ndi.video/tools/",
			sender.Name(), w, h, fps)
	}
	return nil
}

// Enabled reports whether the sender is live.
func (m *Manager) Enabled() bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.started
}

// Shutdown closes the sender and tears down the NDI library.
func (m *Manager) Shutdown() {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.sender != nil {
		m.sender.Close()
		m.sender = nil
	}
	m.started = false
	Shutdown()
}

// PushFrameBGRA reads the current raylib framebuffer (after EndDrawing) and
// submits it to NDI. It converts the BGRA readback to UYVY (4:2:2) — NDI
// then compresses and transports the frame itself. Called from the render
// loop, on the raylib/GL thread.
func (m *Manager) PushFrameBGRA() {
	if m == nil || !m.Enabled() {
		return
	}
	w := int32(rl.GetScreenWidth())
	h := int32(rl.GetScreenHeight())
	if w <= 0 || h <= 0 {
		return
	}

	img := rl.LoadImageFromScreen()
	if img == nil || img.Data == nil || img.Width <= 0 || img.Height <= 0 {
		return
	}
	srcW, srcH := int32(img.Width), int32(img.Height)
	// NDI UYVY requires an even width; the converter crops odd widths.
	srcW = even(srcW)

	m.mu.Lock()
	outW, outH := m.outW, m.outH
	m.mu.Unlock()

	// Scale the framebuffer to the requested NDI output resolution when one
	// is configured (windowed --ndi-width/--ndi-height). Scaling happens in
	// raylib (GPU, bicubic) before the CPU-side UYVY conversion.
	if outW > 0 && outH > 0 && (outW != srcW || outH != srcH) {
		scaled := rl.ImageCopy(img)
		if scaled.Data == nil {
			return
		}
		rl.ImageResize(scaled, outW, outH)
		if scaled.Data == nil || scaled.Width <= 0 || scaled.Height <= 0 {
			rl.UnloadImage(scaled)
			return
		}
		rl.UnloadImage(img)
		img = scaled
		srcW, srcH = even(int32(img.Width)), int32(img.Height)
	}

	// LoadImageFromScreen returns R8G8B8A8 format pixels (see raylib
	// LoadImageFromScreen: format = PIXELFORMAT_UNCOMPRESSED_R8G8B8A8).
	// glReadPixels is called with GL_RGBA, which on little-endian GPUs
	// returns bytes in R,G,B,A order — so we treat the buffer as RGBA.

	// Lazily allocate the UYVY buffer for the current size.
	m.mu.Lock()
	if m.uyvyBuf == nil || m.curW != srcW || m.curH != srcH {
		m.uyvyBuf = make([]byte, srcW*srcH*2)
		m.curW, m.curH = srcW, srcH
	}
	uyvy := m.uyvyBuf
	m.mu.Unlock()

	ConvertRGBAtoUYVY(uyvy, srcW*2, unsafe.Slice((*byte)(img.Data), int(img.Height)*int(img.Width)*4), img.Width*4, srcW, srcH)
	rl.UnloadImage(img)

	fpsN, fpsD := m.fpsN, m.fpsD
	if fpsN == 0 {
		fpsN, fpsD = 60000, 1001
	}
	m.mu.Lock()
	s := m.sender
	m.mu.Unlock()
	if s == nil {
		return
	}
	s.SendVideo(srcW, srcH, FourCCTypeUYVY, srcW*2, fpsN, fpsD, unsafe.Pointer(&uyvy[0]))

	m.frames++
	now := time.Now()
	if m.lastLog.IsZero() {
		m.lastLog = now
	} else if now.Sub(m.lastLog) > 30*time.Second {
		log.Printf("NDI®: %d frames sent, %d receivers connected (%s)",
			m.frames, s.ConnectionCount(), s.Name())
		m.lastLog = now
	}
}

// fpsToFraction converts an integer FPS to N/D (60 -> 60000/1001 is common,
// but plain 60/1 is fine for NDI).
func fpsToFraction(fps int32) (int32, int32) {
	if fps <= 0 {
		return 60000, 1001
	}
	switch fps {
	case 24:
		return 24000, 1001
	case 30:
		return 30000, 1001
	case 60:
		return 60000, 1001
	case 120:
		return 120000, 1001
	}
	return fps, 1
}
