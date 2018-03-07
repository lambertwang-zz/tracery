package main

import "math"

type pointLight struct {
	center vector
	radius float64
	color  floatColor
}

func (l pointLight) light(incident vector, normal vector, s scene) (ray, float64) {
	// Compute intersections
	toRay := lineToRay(incident, l.center)
	nearestT, _ := nearestIntersection(toRay, s)

	v := subtractVector(l.center, incident)

	if v.magnitude() > subtractVector(l.center, toRay.incident(nearestT)).magnitude() {
		return toRay, 0.0
	}

	distance := subtractVector(incident, l.center).magnitude()
	attenuation := 1.0 / (1.0 + (2.0/l.radius)*distance + (1.0/math.Pow(l.radius, 2))*math.Pow(distance, 2))

	return toRay, attenuation
}

func (l pointLight) getColor() floatColor {
	return l.color
}
