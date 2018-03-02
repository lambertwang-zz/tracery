package main

import (
	"image/color"
	"math"
	"strconv"
)

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func vecStr(v vector) string {
	return "x = " + floatStr(v.x) +
		", y = " + floatStr(v.y) +
		", z = " + floatStr(v.z)
}

func meanColor(colors []color.RGBA) color.RGBA {
	var r, g, b, a int

	for _, c := range colors {
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
		a += int(c.A)
	}
	r /= len(colors)
	g /= len(colors)
	b /= len(colors)
	a /= len(colors)

	return color.RGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		uint8(a),
	}
}

func scaleColor(c color.RGBA, f float64) color.RGBA {
	scale := math.Max(0.0, f)
	return color.RGBA{
		uint8(math.Min(float64(c.R)*scale, 255)),
		uint8(math.Min(float64(c.G)*scale, 255)),
		uint8(math.Min(float64(c.B)*scale, 255)),
		c.A,
	}
}

func clamp16to8(x uint16) uint8 {
	if x > 255 {
		return 255
	}
	return uint8(x)
}

func addColor(lhs color.RGBA, rhs color.RGBA) color.RGBA {
	r := clamp16to8(uint16(lhs.R) + uint16(rhs.R))
	g := clamp16to8(uint16(lhs.G) + uint16(rhs.G))
	b := clamp16to8(uint16(lhs.B) + uint16(rhs.B))
	a := clamp16to8(uint16(lhs.A) + uint16(rhs.A))
	return color.RGBA{r, g, b, a}
}
