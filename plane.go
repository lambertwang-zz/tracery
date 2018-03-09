package main

import (
	"math"
)

type plane struct {
	material
	normal vector
	dist   float64
}

func planeIntersect(normal vector, dist float64, ray ray) (t float64) {
	denom := dotProduct(ray.dir, normal)
	if denom >= 0 {
		return -1
	}
	t = -(dotProduct(ray.origin, normal) + dist) / denom
	return
}

func (p plane) intersect(ray ray) float64 {
	return planeIntersect(p.normal, p.dist, ray)
}

func (p plane) traceTo(t float64, params traceParams) (vector, vector) {
	return params.rayCast.incident(t), p.normal
}

func (p plane) sampleC(incident, normal, dir vector, sc scene) floatColor {
	outColor := p.color
	_, modx := math.Modf(incident.x + 100)
	_, modz := math.Modf(incident.z + 100)

	// divx, modx := math.Modf(incident.x + 100)
	// divy, modz := math.Modf(incident.z + 100)
	if (modx-0.5)*(modz-0.5) > 0 {
		/*
			_, red := math.Modf((128 + 16*divx) / 256)
			_, green := math.Modf((128 + 16*divy) / 256)
			outColor = floatColor{
				red * 256,
				green * 256,
				0, 1.0,
			}
		*/
		outColor = floatColor{
			128,
			128,
			128,
			1.0,
		}
	} else {
		outColor = floatColor{
			255,
			255,
			255,
			1.0,
		}
	}

	return outColor
}

func (p plane) bounds(sceneBounds aabb) aabb {
	return sceneBounds
}
