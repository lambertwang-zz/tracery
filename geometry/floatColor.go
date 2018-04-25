package geometry

import (
	"image/color"
	"math"
)

// FloatColor is a color using floating point values
type FloatColor struct {
	R, G, B, A float64
}

// MeanColor TODO: Comment
func MeanColor(colors []FloatColor) FloatColor {
	var r, g, b, a float64

	for _, c := range colors {
		r += c.R
		g += c.G
		b += c.B
		a += c.A
	}
	r /= float64(len(colors))
	g /= float64(len(colors))
	b /= float64(len(colors))
	a /= float64(len(colors))

	return FloatColor{
		r,
		g,
		b,
		a,
	}
}

// RgbaToFloat TODO: Comment
func RgbaToFloat(a color.RGBA) FloatColor {
	return FloatColor{
		float64(a.R),
		float64(a.G),
		float64(a.B),
		float64(a.A) / 255.0,
	}
}

// Scale TODO: Comment
func (a FloatColor) Scale(f ...float64) FloatColor {
	scale := 1.0
	for _, val := range f {
		scale *= val
	}
	scale = math.Max(0.0, scale)
	return FloatColor{
		a.R * scale,
		a.G * scale,
		a.B * scale,
		a.A,
	}
}

// Add TODO: Comment
func (a FloatColor) Add(b FloatColor) FloatColor {
	// Composite a over b
	/*
		inv := b.a * (1 - a.a)
		denom := 1 / (a.a + inv)
		return FloatColor{
			(a.r*a.a + b.r*inv) * denom,
			(a.g*a.a + b.g*inv) * denom,
			(a.b*a.a + b.b*inv) * denom,
			denom,
		}
	*/
	return FloatColor{
		(a.R*a.A + b.R*b.A),
		(a.G*a.A + b.G*b.A),
		(a.B*a.A + b.B*b.A),
		math.Max(a.A, b.A),
	}
}

// ToRgba TODO: Comment
func (a FloatColor) ToRgba() color.RGBA {
	return color.RGBA{
		uint8(math.Max(math.Min(255, a.R), 0)),
		uint8(math.Max(math.Min(255, a.G), 0)),
		uint8(math.Max(math.Min(255, a.B), 0)),
		uint8(math.Max(math.Min(255, a.A*255), 0)),
	}
}
