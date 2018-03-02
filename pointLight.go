package main

import "math"

type pointLight struct {
	center vector
	radius float64
}

func (l pointLight) light(x ray, normal vector, s scene) (ray, float64) {
	// Compute intersections
	nearestT := math.MaxFloat64
	toRay := lineToRay(x.origin, l.center)
	for _, shape := range s.shapes {
		t0, t1 := shape.intersect(toRay)
		// Assume ray cannot cast from inside of the sphere
		if t0 > 0.01 && t1 > 0.01 {
			closeT := math.Min(t0, t1)
			if closeT < nearestT {
				nearestT = closeT
			}
		}
	}

	v := subtractVector(l.center, x.origin)

	if v.magnitude() > subtractVector(l.center, toRay.incident(nearestT)).magnitude() {
		return toRay, 0.0
	}

	distance := subtractVector(x.origin, l.center).magnitude()
	attenuation := 1.0 / (1.0 + (2.0/l.radius)*distance + (1.0/math.Pow(l.radius, 2))*math.Pow(distance, 2))

	return toRay, attenuation
}
