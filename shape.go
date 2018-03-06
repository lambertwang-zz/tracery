package main

import "image/color"

type shape interface {
	shouldTest(ray) bool
	intersect(ray) float64
	// Returns the incident point, all reflected directions, the color of the
	// object at that point, the reflectance, and the normal
	reflect(float64, traceParams, scene) (vector, []traceParams, floatColor, vector)
	getMaterial() material
}

type material struct {
	color       floatColor
	reflectance float64

	// Shading properties
	diffuse  float64
	specular float64
	hardness float64

	refractionIndex float64
}

func defaultMaterial() material {
	return material{floatColor{192, 192, 192, 1.0}, 1.0, 1.0, 1.0, 32, 1.0}
}

func defaultShader(c color.RGBA, r float64) material {
	return material{rgbaToFloat(c), r, 1.0, 1.0, 32, 1.0}
}

func createMaterial(c color.RGBA, r float64, diff float64, spec float64, hard float64, refr float64) material {
	return material{rgbaToFloat(c), r, diff, spec, hard, refr}
}
