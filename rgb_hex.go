package pickem

import "fmt"

// RGBHex represents a typical 24-bit RGB color in hexidecimal format.
type RGBHex string

func hex(r, g, b int8) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func colors(hex string) (r, g, b int8, err error) {
	_, err = fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b)
	return
}

// MakeRGBHex creates an RGBHex object from separate r, g, and b channels.
func MakeRGBHex(r, g, b int8) RGBHex {
	return RGBHex(hex(r, g, b))
}

// RGBA implements the color.Color interface.
func (c RGBHex) RGBA() (r, g, b, a uint32) {
	rr, gg, bb, _ := colors(string(c))
	r = uint32(rr)
	r |= r << 8
	g = uint32(gg)
	g |= g << 8
	b = uint32(bb)
	b |= b << 8
	a = uint32(0xffff)
	return
}
