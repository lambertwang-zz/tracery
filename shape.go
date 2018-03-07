package main

import (
	"image/color"
	"math"
)

type shape interface {
	getMaterial() material

	shouldTest(ray) bool
	intersect(ray) float64
	// Returns the incident and the normal
	traceTo(t float64, params traceParams) (vector, vector)
	sampleC(incident, normal, dir vector, sc scene) floatColor
	// Returns the incident point, all reflected directions, the color of the
	// object at that point, the reflectance, and the normal
	reflect(incident, normal, dir vector, sc scene) vector
	refract(incident, normal, dir vector, sc scene) vector
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
	return material{floatColor{192, 192, 192, 1.0}, 0.0, 1.0, 1.0, 32, 0.0}
}

func defaultShader(c color.RGBA, r float64) material {
	return material{rgbaToFloat(c), r, 1.0, 1.0, 32, 0.0}
}

func createMaterial(c color.RGBA, r float64, diff float64, spec float64, hard float64, refr float64) material {
	return material{rgbaToFloat(c), r, diff, spec, hard, refr}
}

func (m material) getMaterial() material {
	return m
}

func (m material) sampleC(incident, normal, dir vector, sc scene) floatColor {
	return m.color
}

func (m material) reflect(incident, normal, dir vector, sc scene) (reflection vector) {
	reflection = subtractVector(dir, normal.scale(2*dotProduct(normal, dir)))
	return
}

func (m material) refract(incident, normal, dir vector, sc scene) (refraction vector) {
	refractionIndex := 1.0 / m.refractionIndex
	/*
		theta1 := math.Acos(dotProduct(dir, normal))
		theta2 := math.Asin(math.Sin(theta1) * refractionIndex)
		refraction := subtractVector(
			addVectors(dir.scale(math.Sin(theta2)/math.Sin(theta1)), normal.scale(math.Cos(theta1)*math.Sin(theta2)/math.Sin(theta1))),
			normal.scale(math.Cos(theta2)),
		).norm()
	*/

	c1 := dotProduct(normal, dir)
	c2 := math.Sqrt(1 - math.Pow(refractionIndex, 2)*(1-math.Pow(c1, 2)))
	refraction = addVectors(dir.scale(refractionIndex), normal.scale(refractionIndex*c1-c2)).norm()
	return
}
