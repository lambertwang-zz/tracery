package main

import (
	"image/color"
	"math"
)

const reflectanceThreshold = 0.01

type traceParams struct {
	rayCast ray
	depth   int
	value   float64
}

func traceSample(sa sample, t traceParams, s scene) (color.RGBA, float64) {
	c := make([]floatColor, len(sa.dispatchedRays))

	outDepth := -1.0
	for i, r := range sa.dispatchedRays {
		newParams := traceParams{
			r, t.depth, t.value,
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
	if t.value <= 0.001 {
		return floatColor{0, 0, 0, 0}, math.MaxFloat64
	}
	nearestT, index := nearestIntersection(t.rayCast, s)
	if index < 0 {
		return floatColor{0, 0, 0, 1.0}, nearestT
	}
	if t.depth <= 0 {
		return floatColor{0, 0, 0, 1.0}, nearestT
	}

	incident, normal := s.shapes[index].traceTo(nearestT, t)

	m := s.shapes[index].getMaterial()

	reflectionColor := floatColor{0, 0, 0, 0}
	reflectionDir := s.shapes[index].reflect(incident, normal, t.rayCast.dir, s)
	if m.reflectance > 0.001 {
		reflectionTrace := traceParams{
			ray{incident, reflectionDir},
			t.depth - 1,
			t.value * m.reflectance,
		}
		reflectionColor, _ = trace(reflectionTrace, s)
	}
	refractionColor := floatColor{0, 0, 0, 0}
	if m.refractionIndex > 0.001 {
		refractionDir := s.shapes[index].refract(incident, normal, t.rayCast.dir, s)
		refractionTrace := traceParams{
			ray{incident, refractionDir},
			t.depth - 1,
			t.value,
		}
		refractionColor, _ = trace(refractionTrace, s)
	}
	phongColor := s.shapes[index].sampleC(incident, normal, t.rayCast.dir, s)

	outColor := phongColor.scale(s.ambientLight)
	for _, light := range s.lights {
		toRay, attenuation := light.light(incident, normal, s)

		if m.diffuse > 0.0 {
			diffuse := m.diffuse * attenuation * math.Max(0, dotProduct(normal, toRay.dir))
			outColor = outColor.add(phongColor.scale(diffuse))
		}
		// Blinn-Phong shading
		if m.specular > 0.0 {
			// Calculate the light reflection for specular lighting
			halfway := addVectors(t.rayCast.dir.neg(), toRay.dir).norm()
			specular := m.specular * attenuation * math.Pow(math.Max(0.0, dotProduct(halfway, normal)), m.hardness)
			// specular := m.specular * attenuation * math.Pow(math.Max(0.0, dotProduct(toRay.dir, t.rayCast.dir)), m.hardness)
			// specular = 0
			outColor = outColor.add(light.getColor().scale(specular))
		}
	}

	outColor = outColor.add(reflectionColor).add(refractionColor).scale(t.value)
	// outColor = phongColor.scale(t.value, lightVal)
	return outColor, nearestT
}
