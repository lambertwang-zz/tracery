package main

import (
	"image/color"
	"math"
)

const reflectanceThreshold = 0.01

type traceParams struct {
	reflections int
	reflectance float64
	color       color.RGBA
}

func traceSample(sa sample, t traceParams, s scene) (color.RGBA, float64) {
	c := make([]color.RGBA, len(sa.casts))
	outDepth := -1.0
	for i, r := range sa.casts {
		c[i], outDepth = trace(r, t, s)
	}
	return meanColor(c), outDepth
}

func trace(ray ray, t traceParams, s scene) (color.RGBA, float64) {
	// Compute intersections
	nearestT := math.MaxFloat64
	index := -1
	for i, shape := range s.shapes {
		t0, t1 := shape.intersect(ray)
		// Assume ray cannot cast from inside of the sphere
		if t0 > 0.01 && t1 > 0.01 {
			closeT := math.Min(t0, t1)
			if closeT < nearestT {
				nearestT = closeT
				index = i
			}
		}
	}
	if index < 0 {
		return t.color, nearestT
	}

	if t.reflections <= 0 {
		return t.color, nearestT
	}

	reflection, c, reflectance, normal := s.shapes[index].reflect(ray, nearestT, s)
	lightVal := s.ambientLight
	m := s.shapes[index].getMaterial()
	for _, light := range s.lights {
		toRay, attenuation := light.light(reflection, normal, s)

		diffuse := 0.0
		if m.diffuse > 0.0 {
			diffuse = m.diffuse * attenuation * math.Max(0, dotProduct(normal, toRay.dir))
		}
		// Blinn-Phon shading
		specular := 0.0
		if m.specular > 0.0 {
			// Calculate the light reflection for specular lighting
			halfway := addVectors(ray.dir.neg(), toRay.dir).norm()
			specular = m.specular * attenuation * math.Pow(math.Max(0.0, dotProduct(halfway, normal)), m.hardness)
		}

		lightVal += specular + diffuse
	}

	t.reflections--
	t.color = addColor(t.color, scaleColor(scaleColor(c, t.reflectance), lightVal))
	t.reflectance *= reflectance
	if reflectance < reflectanceThreshold {
		return t.color, nearestT
	}

	outColor, _ := trace(reflection, t, s)
	return outColor, nearestT
}
