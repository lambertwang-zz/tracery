package main

import (
	"image/color"
	"math"
)

type sphere struct {
	material
	center vector
	radius float64
}

func (s sphere) intersect(ray ray) (float64, float64) {
	v := subtractVector(ray.origin, s.center)
	vd := dotProduct(v, ray.dir)
	discriminant := math.Pow(vd, 2) - (math.Pow(v.magnitude(), 2) - math.Pow(s.radius, 2))

	// If the discriminant is negative, the ray does not intersect the sphere
	if discriminant < 0 {
		return -1, -1
	}

	sqrtd := math.Sqrt(discriminant)
	return -vd + sqrtd, -vd - sqrtd
}

func (s sphere) reflect(x ray, t float64, sc scene) (ray, color.RGBA, float64) {
	// n = (y - c) / || y - c|| where y = s + td
	incident := subtractVector(addVectors(x.origin, x.dir.scale(t)), s.center)
	normal := incident.norm()
	reflection := subtractVector(x.dir, normal.scale(2*(dotProduct(normal, x.dir))))
	outRay := ray{addVectors(incident, s.center), reflection}

	var lightVal float64
	for _, light := range sc.lights {
		lightVal += light.light(outRay, normal, sc)
	}

	return outRay, scaleColor(s.color, lightVal), s.reflectance
}
