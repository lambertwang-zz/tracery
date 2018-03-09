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

func nearestIntersection(r ray, s scene) (nearestT float64, nearestShape *shape) {
	nearestT, nearestShape = checkBvh(r, s.bvh)
	// Following code is before BVH
	/*
		nearestT = math.MaxFloat64
		nearestShape = nil
		for i, shape := range s.shapes {
			t := shape.intersect(r)
			if t > 0.001 && t < nearestT {
				nearestT = t
				nearestShape = &s.shapes[i]
			}
		}
	*/
	return
}

// Returns the nearestT and the shape intersected
func checkBvh(r ray, n *bvhNode) (float64, *shape) {
	intersectsBounds, _, _ := n.bounds.intersects(r)
	if intersectsBounds {
		if n.s == nil {
			lhsT, lhsS := checkBvh(r, n.lhs)
			rhsT, rhsS := checkBvh(r, n.rhs)
			if lhsS != nil && lhsT < rhsT {
				return lhsT, lhsS
			}
			if rhsS != nil {
				return rhsT, rhsS
			}
		} else {
			t := (*n.s).intersect(r)
			if t > 0.001 {
				return t, n.s
			}
		}
	}
	return math.MaxFloat64, nil
}

func trace(t traceParams, s scene) (floatColor, float64) {
	if t.value <= 0.001 {
		return floatColor{0, 0, 0, 0}, math.MaxFloat64
	}
	nearestT, nearestShape := nearestIntersection(t.rayCast, s)
	if nearestShape == nil {
		return floatColor{0, 0, 0, 1.0}, nearestT
	}
	if t.depth <= 0 {
		return floatColor{0, 0, 0, 1.0}, nearestT
	}

	incident, normal := (*nearestShape).traceTo(nearestT, t)

	m := (*nearestShape).getMaterial()

	reflectionColor := floatColor{0, 0, 0, 0}
	reflectionDir := (*nearestShape).reflect(incident, normal, t.rayCast.dir, s)
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
		refractionDir := (*nearestShape).refract(incident, normal, t.rayCast.dir, s)
		refractionTrace := traceParams{
			ray{incident, refractionDir},
			t.depth - 1,
			t.value,
		}
		refractionColor, _ = trace(refractionTrace, s)
	}
	phongColor := (*nearestShape).sampleC(incident, normal, t.rayCast.dir, s)

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
