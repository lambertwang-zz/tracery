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

func (s sphere) reflect(t float64, params traceParams, sc scene) (vector, []traceParams, floatColor, vector) {
	// n = (y - c) / || y - c|| where y = s + td
	x := params.rayCast
	incident := x.incident(t)
	normal := subtractVector(incident, s.center).norm()
	c1 := dotProduct(normal, x.dir)
	reflection := subtractVector(x.dir, normal.scale(2*c1))
	var reflectedDirs []traceParams
	reflectedDirs = append(reflectedDirs, traceParams{
		ray{incident, reflection},
		params.reflections - 1,
		params.reflectance * s.reflectance,
		params.refractance,
	})
	if s.color.a < 1.0 {
		refractionIndex := params.refractance / s.refractionIndex
		theta1 := math.Acos(dotProduct(x.dir, normal))
		theta2 := math.Asin(math.Sin(theta1) * refractionIndex)
		refraction := subtractVector(
			addVectors(x.dir.scale(math.Sin(theta2)/math.Sin(theta1)), normal.scale(math.Cos(theta1)*math.Sin(theta2)/math.Sin(theta1))),
			normal.scale(math.Cos(theta2)),
		).norm()
		// c2 := math.Sqrt(1 - math.Pow(refractionIndex, 2)*(1-math.Pow(c1, 2)))
		// refraction := addVectors(x.dir.scale(refractionIndex), normal.scale(refractionIndex*c1-c2)).norm()
		newRefractionIndex := s.refractionIndex
		if c1 < 0 {
			newRefractionIndex = 1.0
		}
		reflectedDirs = append(reflectedDirs, traceParams{
			ray{incident, refraction},
			params.reflections - 1,
			params.reflectance,
			newRefractionIndex,
		})
	}
	return incident, reflectedDirs, s.color.scale(params.reflectance), normal
}

func (s sphere) getMaterial() material {
	return s.material
}
