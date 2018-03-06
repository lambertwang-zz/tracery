package main

import (
	"image/color"
	"math"
)

const reflectanceThreshold = 0.01

type traceParams struct {
	rayCast     ray
	reflections int
	reflectance float64
	refractance float64
}

func traceSample(sa sample, t traceParams, s scene) (color.RGBA, float64) {
	c := make([]floatColor, len(sa.casts))

	outDepth := -1.0
	for i, r := range sa.casts {
		newParams := traceParams{
			r, t.reflections, t.reflectance, t.refractance,
		}
		c[i], outDepth = trace(newParams, s)
	}
	return meanColor(c).toRgba(), outDepth
}

func nearestIntersection(r ray, s scene) (nearestT float64, index int) {
	nearestT = math.MaxFloat64
	index = -1
	for i, shape := range s.shapes {
		t := shape.intersect(r)
		// Use some minimum distance the ray must travel before it can intersect
		if t > 0.001 {
			if t < nearestT {
				nearestT = t
				index = i
			}
		}
	}
	return
}

func trace(t traceParams, s scene) (floatColor, float64) {
	if t.reflectance <= 0.001 {
		return floatColor{0, 0, 0, 0}, math.MaxFloat64
	}
	nearestT, index := nearestIntersection(t.rayCast, s)
	if index < 0 {
		return floatColor{0, 0, 0, 0}, nearestT
	}
	if t.reflections <= 0 {
		return floatColor{0, 0, 0, 0}, nearestT
	}

	incident, reflectionParams, c, normal := s.shapes[index].reflect(nearestT, t, s)
	lightVal := s.ambientLight
	m := s.shapes[index].getMaterial()
	for _, light := range s.lights {
		toRay, attenuation := light.light(incident, normal, s)

		diffuse := 0.0
		if m.diffuse > 0.0 {
			diffuse = m.diffuse * attenuation * math.Max(0, dotProduct(normal, toRay.dir))
		}
		// Blinn-Phon shading
		specular := 0.0
		if m.specular > 0.0 {
			// Calculate the light reflection for specular lighting
			halfway := addVectors(t.rayCast.dir.neg(), toRay.dir).norm()
			specular = m.specular * attenuation * math.Pow(math.Max(0.0, dotProduct(halfway, normal)), m.hardness)
		}

		lightVal += specular + diffuse
	}

	var outColors []floatColor
	for _, newParam := range reflectionParams {
		newColor, _ := trace(newParam, s)
		outColors = append(outColors, newColor)
	}
	// outColor := c.scale(t.reflectance, lightVal).add(meanColor(outColors))
	outColor := meanColor(outColors).add(c.scale(t.reflectance, lightVal))

	return outColor, nearestT
}
