package ndi

import "testing"

// TestConvertBasicColors checks a small 2x2 frame with primary colors.
// GL readback stores rows bottom-up: the first stored row is the image's
// bottom row and the converter's flip makes it dst row yres-1.
//   stored row 0 (image bottom): blue (0,0,255)   white (255,255,255)
//   stored row 1 (image top):    red (255,0,0)    green (0,255,0)
// So dst row 0 (top) holds red/green and dst row 1 (bottom) blue/white.
func TestConvertBasicColors(t *testing.T) {
	src := make([]byte, 2*2*4)
	put := func(px int, r, g, b byte) {
		src[px*4] = r
		src[px*4+1] = g
		src[px*4+2] = b
		src[px*4+3] = 255
	}
	// stored row 0 = bottom row: blue, white
	put(0, 0, 0, 255)
	put(1, 255, 255, 255)
	// stored row 1 = top row: red, green
	put(2, 255, 0, 0)
	put(3, 0, 255, 0)

	dst := make([]byte, 2*2*2)
	ConvertRGBAtoUYVY(dst, 4, src, 8, 2, 2)

	// Y from BT.601 full-range coefficients (x*255+128)>>8:
	// red -> 77, green -> 149, blue -> 29, white -> 255.
	if dst[1] != 77 {
		t.Errorf("row0 px0: Y = %d, want 77 (red)", dst[1])
	}
	if dst[3] != 149 {
		t.Errorf("row0 px1: Y = %d, want 149 (green)", dst[3])
	}
	if dst[4+1] != 29 {
		t.Errorf("row1 px0: Y = %d, want 29 (blue)", dst[4+1])
	}
	if dst[4+3] != 255 {
		t.Errorf("row1 px1: Y = %d, want 255 (white)", dst[4+3])
	}
}

// TestConvertOddWidthCropped ensures odd widths are cropped to even without
// panicking (UYVY needs even width).
func TestConvertOddWidthCropped(t *testing.T) {
	xres, yres := int32(3), int32(2)
	src := make([]byte, xres*yres*4)
	dst := make([]byte, 2*yres*2) // even(3)=2 -> 2 px per row
	ConvertRGBAtoUYVY(dst, 4, src, xres*4, xres, yres)
}

// TestConvertChroma checks the shared chroma of a red/green pair.
// Avg RGB = (128, 128, 0):
//   U = (-43*128 - 85*128)/256 = -64 -> stored 64
//   V = (128*128 - 107*128)/256 = 10 -> stored 138
func TestConvertChroma(t *testing.T) {
	src := []byte{
		255, 0, 0, 255, 0, 255, 0, 255,
	}
	dst := make([]byte, 4)
	ConvertRGBAtoUYVY(dst, 4, src, 8, 2, 1)
	if got := dst[0]; got != 64 {
		t.Errorf("U = %d, want 64", got)
	}
	if got := dst[2]; got != 138 {
		t.Errorf("V = %d, want 138", got)
	}
}
