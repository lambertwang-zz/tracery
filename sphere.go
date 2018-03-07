package main

import (
	"math"
)

type sphere struct {
	material
	center vector
	radius float64
}

func (s sphere) shouldTest(ray ray) bool {
	return dotProduct(subtractVector(s.center, ray.origin), ray.dir) > 0
}

func (s sphere) intersect(ray ray) float64 {
	v := subtractVector(ray.origin, s.center)
	vd := dotProduct(v, ray.dir)
	discriminant := math.Pow(vd, 2) - (math.Pow(v.magnitude(), 2) - math.Pow(s.radius, 2))

	// If the discriminant is negative, the ray does not intersect the sphere
	if discriminant < 0 {
		return -1
	}

	sqrtd := math.Sqrt(discriminant)
	// we are only concerned with the nearest point
	// return -vd + sqrtd, -vd - sqrtd
	return -vd - sqrtd
}

func (s sphere) traceTo(t float64, params traceParams) (incident vector, normal vector) {
	// n = (y - c) / || y - c|| where y = s + td
	incident = params.rayCast.incident(t)
	normal = subtractVector(incident, s.center).norm()
	return
}
