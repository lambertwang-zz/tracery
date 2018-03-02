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

func (p plane) shouldTest(ray ray) bool {
	return true
}

func (p plane) intersect(ray ray) (float64, float64) {
	denom := dotProduct(ray.dir, p.normal)
	if denom >= 0 {
		return -1, -1
	}
	return -(dotProduct(ray.origin, p.normal) + p.dist) / denom, math.MaxFloat64
}

func (p plane) reflect(x ray, t float64, sc scene) (ray, color.RGBA, float64, vector) {
	incident := addVectors(x.origin, x.dir.scale(t))
	outColor := p.color
	_, modx := math.Modf(incident.x + 100)
	_, modz := math.Modf(incident.z + 100)
	if (modx-0.5)*(modz-0.5) > 0 {
		outColor = scaleColor(p.color, .6)
	}
	reflection := subtractVector(x.dir, p.normal.scale(2*(dotProduct(p.normal, x.dir))))

	if incident.z > 20 {
		outColor = color.RGBA{0, 0, 0, 255}
	} else if incident.z > 10 {
		outColor = scaleColor(outColor, (20-incident.z)/10)
	}

	return ray{incident, reflection}, outColor, p.reflectance, p.normal
}

func (p plane) getMaterial() material {
	return p.material
}
