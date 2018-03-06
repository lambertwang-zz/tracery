package main

import (
	"image/color"
	"math"
	"strconv"
)

type floatColor struct {
	r, g, b, a float64
}

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func vecStr(v vector) string {
	return "x = " + floatStr(v.x) +
		", y = " + floatStr(v.y) +
		", z = " + floatStr(v.z)
}

func meanColor(colors []floatColor) floatColor {
	var r, g, b, a float64

	for _, c := range colors {
		r += c.r
		g += c.g
		b += c.b
		a += c.a
	}
	r /= float64(len(colors))
	g /= float64(len(colors))
	b /= float64(len(colors))
	a /= float64(len(colors))

	return floatColor{
		r,
		g,
		b,
		a,
	}
}

func (a floatColor) scale(f ...float64) floatColor {
	scale := 1.0
	for _, val := range f {
		scale *= val
	}
	scale = math.Max(0.0, scale)
	return floatColor{
		a.r * scale,
		a.g * scale,
		a.b * scale,
		a.a,
	}
}

func clamp16to8(x uint16) uint8 {
	if x > 255 {
		return 255
	}
	return uint8(x)
}

func (a floatColor) add(b floatColor) floatColor {
	// Composite a over b
	/*
		inv := b.a * (1 - a.a)
		denom := 1 / (a.a + inv)
		return floatColor{
			(a.r*a.a + b.r*inv) * denom,
			(a.g*a.a + b.g*inv) * denom,
			(a.b*a.a + b.b*inv) * denom,
			denom,
		}
	*/
	return floatColor{
		(a.r*a.a + b.r*b.a),
		(a.g*a.a + b.g*b.a),
		(a.b*a.a + b.b*b.a),
		math.Max(a.a, b.a),
	}
}

func rgbaToFloat(a color.RGBA) floatColor {
	return floatColor{
		float64(a.R),
		float64(a.G),
		float64(a.B),
		float64(a.A) / 255.0,
	}
}

func (a floatColor) toRgba() color.RGBA {
	return color.RGBA{
		uint8(math.Max(math.Min(255, a.r), 0)),
		uint8(math.Max(math.Min(255, a.g), 0)),
		uint8(math.Max(math.Min(255, a.b), 0)),
		uint8(math.Max(math.Min(255, a.a*255), 0)),
	}
}
