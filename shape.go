package main

import "image/color"

type shape interface {
	shouldTest(ray) bool
	intersect(ray) (float64, float64)
	reflect(ray, float64, scene) (ray, color.RGBA, float64, vector)
	getMaterial() material
}

type material struct {
	color       color.RGBA
	reflectance float64

	// Shading properties
	diffuse  float64
	specular float64
	hardness float64
}

func defaultMaterial() material {
	return material{color.RGBA{192, 192, 192, 255}, 1.0, 1.0, 1.0, 32}
}

func defaultShader(c color.RGBA, r float64) material {
	return material{c, r, 1.0, 1.0, 32}
}

func createMaterial(c color.RGBA, r float64, diff float64, spec float64, hard float64) material {
	return material{c, r, diff, spec, hard}
}
