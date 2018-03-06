package main

import (
	"math"
)

type plane struct {
	material
	normal vector
	dist   float64
}

func (p plane) shouldTest(ray ray) bool {
	return true
}

func planeIntersect(normal vector, dist float64, ray ray) float64 {
	denom := dotProduct(ray.dir, normal)
	if denom >= 0 {
		return -1
	}
	t := -(dotProduct(ray.origin, normal) + dist) / denom
	incident := ray.incident(t)
	if incident.z > 20 {
		return -1
	}
	return t
}

func (p plane) intersect(ray ray) float64 {
	return planeIntersect(p.normal, p.dist, ray)
}

func (p plane) reflect(t float64, params traceParams, sc scene) (vector, []traceParams, floatColor, vector) {
	x := params.rayCast
	incident := x.incident(t)
	outColor := p.color
	divx, modx := math.Modf(incident.x + 100)
	divy, modz := math.Modf(incident.z + 100)
	if (modx-0.5)*(modz-0.5) > 0 {
		_, red := math.Modf((128 + 16*divx) / 256)
		_, green := math.Modf((128 + 16*divy) / 256)
		outColor = floatColor{
			red * 256,
			green * 256,
			0, 1.0,
		}
	} else {
		outColor = floatColor{
			255,
			255,
			255,
			1.0,
		}
	}
	reflection := subtractVector(x.dir, p.normal.scale(2*(dotProduct(p.normal, x.dir))))

	if incident.z > 20 {
		outColor = floatColor{255, 255, 255, 1.0}
	} else if incident.z > 10 {
		outColor = outColor.scale((20 - incident.z) / 10)
	}

	return incident, []traceParams{traceParams{
		ray{incident, reflection},
		params.reflections - 1,
		params.reflectance * p.reflectance,
		params.refractance,
	}}, outColor.scale(params.reflectance), p.normal
}

func (p plane) getMaterial() material {
	return p.material
}
