package ndi

// raylib's screen readback (glReadPixels with GL_RGBA) yields one byte per
// channel in R,G,B,A order per pixel (see rlReadScreenPixels in raylib).
// Rows come back bottom-up and the alpha channel is forced to 255 by raylib.
const rgbaBytesPerPixel = 4

// ConvertRGBAtoUYVY converts a top-down RGBA (4 bytes/px) frame to UYVY
// (4:2:2, BT.601) for NDI, removing the OpenGL bottom-up row order on the fly.
//
// UYVY chroma is shared by pixel pairs, so NDI requires an even frame width:
// an odd source width is cropped by one column on the right.
//
// src must contain yres rows of srcStride bytes (raylib readbacks have
// srcStride = srcWidth*4 and no extra padding). dst must be large enough for
// yres rows of dstStride bytes; width must not exceed the space available.
// Chroma is shared by pixel pairs and averaged over both pixels; luma for
// each pixel is computed individually.
// UYVY requires an even pixel width; NDI rejects odd widths. Crop one column.
func even(x int32) int32 {
	if x%2 != 0 {
		return x - 1
	}
	return x
}

func ConvertRGBAtoUYVY(dst []byte, dstStride int32, src []byte, srcStride int32, xres, yres int32) {
	// Integer BT.601 full range fixed point conversion:
	//   Y = ( 77*R + 150*G +  29*B + 128) >> 8
	//   U = (-43*R -  85*G + 128*B + 128) >> 8   (add 128 to store)
	//   V = (128*R - 107*G -  21*B + 128) >> 8   (add 128 to store)
	const (
		yR = 77
		yG = 150
		yB = 29
		uR = -43
		uG = -85
		uB = 128
		vR = 128
		vG = -107
		vB = -21
	)
	xres = even(xres)

	for y := int32(0); y < yres; y++ {
		// OpenGL framebuffer reads are bottom-up; flip to top-down for NDI.
		srcRow := (yres - 1 - y) * srcStride
		dstRow := y * dstStride
		var x int32
		for x = 0; x+1 < xres; x += 2 {
			i0 := srcRow + x*rgbaBytesPerPixel
			i1 := i0 + rgbaBytesPerPixel

			r0, g0, b0 := int32(src[i0]), int32(src[i0+1]), int32(src[i0+2])
			r1, g1, b1 := int32(src[i1]), int32(src[i1+1]), int32(src[i1+2])

			y0 := (yR*r0 + yG*g0 + yB*b0 + 128) >> 8
			y1 := (yR*r1 + yG*g1 + yB*b1 + 128) >> 8

			rAvg := (r0 + r1) >> 1
			gAvg := (g0 + g1) >> 1
			bAvg := (b0 + b1) >> 1

			u := clampChroma((uR*rAvg + uG*gAvg + uB*bAvg) >> 8)
			v := clampChroma((vR*rAvg + vG*gAvg + vB*bAvg) >> 8)

			di := dstRow + x*2
			// UYVY byte order: U0, Y0, V0, Y1
			dst[di+0] = byte(u + 128)
			dst[di+1] = byte(clampLuma(y0))
			dst[di+2] = byte(v + 128)
			dst[di+3] = byte(clampLuma(y1))
		}
	}
}

func clampLuma(v int32) int32 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func clampChroma(v int32) int32 {
	if v < -128 {
		return -128
	}
	if v > 127 {
		return 127
	}
	return v
}
