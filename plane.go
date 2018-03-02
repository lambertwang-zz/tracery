package main

import (
	"image/color"
	"math"
)

type plane struct {
	material
	normal vector
	dist   float64
}

func (p plane) intersect(ray ray) (float64, float64) {
	denom := dotProduct(ray.dir, p.normal)
	if denom >= 0 {
		return -1, -1
	}
	return -(dotProduct(ray.origin, p.normal) + p.dist) / denom, math.MaxFloat64
}

func (p plane) reflect(x ray, t float64, sc scene) (ray, color.RGBA, float64) {
	incident := addVectors(x.origin, x.dir.scale(t))
	reflection := subtractVector(x.dir, p.normal.scale(2*(dotProduct(p.normal, x.dir))))

	outRay := ray{incident, reflection}

	var lightVal float64
	for _, light := range sc.lights {
		lightVal += light.light(outRay, p.normal, sc)
	}

	return outRay, scaleColor(p.color, lightVal), p.reflectance
}
