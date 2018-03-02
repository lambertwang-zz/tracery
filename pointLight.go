package main

import "math"

type pointLight struct {
	center vector
}

func (l pointLight) light(x ray, normal vector, s scene) float64 {
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
		return 0.0
	}

	return math.Max(0.0, dotProduct(normal, toRay.dir))
}
