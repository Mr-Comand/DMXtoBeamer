package animation_helpers

import (
	"fmt"
	"image/color"
	"math"
	"strconv"
	"strings"
)

type Color struct {
	color.RGBA
}

func AsFullColor[T int | float64](r, g, b T) color.RGBA {
	// Simple color normalization based on max intensity
	peakVolume := math.Max(float64(r), math.Max(float64(g), float64(b)))
	// maxIntensity := math.Max(float64(r), math.Max(float64(g), float64(b)))
	normalizedRed := float64(r) / peakVolume
	normalizedGreen := float64(g) / peakVolume
	normalizedBlue := float64(b) / peakVolume

	return color.RGBA{
		R: uint8(normalizedRed * 255),
		G: uint8(normalizedGreen * 255),
		B: uint8(normalizedBlue * 255),
		A: 255,
	}
}
func AsDynamicColor[T int | float64](r, g, b T, peakVolume float64) (color.RGBA, float64) {
	// Normalize the peak volume to ensure the highest RGB component is used
	peakVolume = math.Max(math.Max(float64(r), math.Max(float64(g), float64(b))), peakVolume)

	// Normalize RGB values based on the peak volume
	normalizedRed := float64(r) / peakVolume
	normalizedGreen := float64(g) / peakVolume
	normalizedBlue := float64(b) / peakVolume

	// Return the RGBA color and updated peak volume
	return color.RGBA{
		R: uint8(normalizedRed * 255),
		G: uint8(normalizedGreen * 255),
		B: uint8(normalizedBlue * 255),
		A: 255,
	}, peakVolume
}

// SetHex sets the color from a hex string (supports #RGB, #RGBA, #RRGGBB, and #RRGGBBAA formats)
func (c *Color) SetHex(hex string) error {
	hex = strings.TrimPrefix(hex, "#")
	switch len(hex) {
	case 3: // #RGB format
		r, _ := strconv.ParseUint(string(hex[0])+string(hex[0]), 16, 8)
		g, _ := strconv.ParseUint(string(hex[1])+string(hex[1]), 16, 8)
		b, _ := strconv.ParseUint(string(hex[2])+string(hex[2]), 16, 8)
		c.RGBA = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	case 4: // #RGB format
		r, _ := strconv.ParseUint(string(hex[0])+string(hex[0]), 16, 8)
		g, _ := strconv.ParseUint(string(hex[1])+string(hex[1]), 16, 8)
		b, _ := strconv.ParseUint(string(hex[2])+string(hex[2]), 16, 8)
		a, _ := strconv.ParseUint(string(hex[3])+string(hex[3]), 16, 8)
		c.RGBA = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	case 6: // #RRGGBB format
		r, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return err
		}
		g, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return err
		}
		b, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return err
		}
		c.RGBA = color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	case 8: // #RRGGBBAA format
		r, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return err
		}
		g, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return err
		}
		b, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return err
		}
		a, err := strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil {
			return err
		}
		c.RGBA = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	default:
		return fmt.Errorf("invalid hex format")
	}
	return nil
}

// Hex returns the color as a hex string (e.g., "#RRGGBB")
func (c Color) Hex() string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}

// FromString creates a Color from a hex string
func FromString(hex string) (Color, error) {
	var c Color
	err := c.SetHex(hex)
	if err != nil {
		return Color{}, err
	}
	return c, nil
}
