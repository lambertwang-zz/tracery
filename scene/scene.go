package scene

import (
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// A Scene is a collection of lights and shapes in the world
// And also the Bounding Volume Hierarchy
type Scene struct {
	shapes       *[]Shape
	lights       *[]Light
	ambientLight float64
	bvh          *bvhNode
}

// TraceParams TODO Comment
type TraceParams struct {
	RayCast g.Ray
	Depth   int
	Value   float64
}

// CreateScene TODO: Comment
func CreateScene(shapes *[]Shape, lights *[]Light, ambient float64, sceneBounds Aabb) (s Scene) {
	s.shapes = shapes
	s.lights = lights
	s.ambientLight = ambient
	s.bvh = constructHeirarchy(shapes, sceneBounds)

	return
}

// Trace TODO: Comment
func (s Scene) Trace(t TraceParams) (outColor g.FloatColor, depth float64) {
	if t.Value <= 0.001 {
		depth = math.MaxFloat64
		return
	}
	depth, nearestShape := NearestIntersection(t.RayCast, s)
	if nearestShape == nil {
		outColor.A = 1.0
		return
	}
	if t.Depth <= 0 {
		outColor.A = 1.0
		return
	}

	incident, normal := (*nearestShape).traceTo(depth, t.RayCast)

	m := (*nearestShape).getMaterial()

	var reflectionColor g.FloatColor
	reflectionDir := (*nearestShape).reflect(incident, normal, t.RayCast.Dir, s)
	if m.reflectance > 0.001 {
		reflectionTrace := TraceParams{
			g.Ray{Origin: incident, Dir: reflectionDir},
			t.Depth - 1,
			t.Value * m.reflectance,
		}
		reflectionColor, _ = s.Trace(reflectionTrace)
	}

	var refractionColor g.FloatColor
	if m.refractionIndex > 0.001 {
		refractionDir := (*nearestShape).refract(incident, normal, t.RayCast.Dir, s)
		refractionTrace := TraceParams{
			g.Ray{Origin: incident, Dir: refractionDir},
			t.Depth - 1,
			t.Value,
		}
		refractionColor, _ = s.Trace(refractionTrace)
	}

	phongColor := (*nearestShape).sampleC(incident, normal, t.RayCast.Dir, s)
	outColor = phongColor.Scale(s.ambientLight)
	for _, light := range *s.lights {
		toRay, attenuation := light.light(incident, normal, s)

		if m.diffuse > 0.0 {
			diffuse := m.diffuse * attenuation * math.Max(0, normal.Dot(toRay.Dir))
			outColor = outColor.Add(phongColor.Scale(diffuse))
		}
		// Blinn-Phong shading
		if m.specular > 0.0 {
			// Calculate the light reflection for specular lighting
			halfway := t.RayCast.Dir.Neg().Plus(toRay.Dir).Norm()
			specular := m.specular * attenuation * math.Pow(math.Max(0.0, halfway.Dot(normal)), m.hardness)
			// specular := m.specular * attenuation * math.Pow(math.Max(0.0, dotProduct(toRay.dir, t.rayCast.dir)), m.hardness)
			// specular = 0
			outColor = outColor.Add(light.GetColor().Scale(specular))
		}
	}

	outColor = outColor.Add(reflectionColor).Add(refractionColor).Scale(t.Value)
	// outColor = phongColor.scale(t.value, lightVal)
	return
}
