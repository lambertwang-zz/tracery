package scene

import (
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// A Plane shape
type Plane struct {
	Material
	Normal g.Vector
	Dist   float64
}

func planeIntersect(normal g.Vector, dist float64, ray g.Ray) (t float64) {
	denom := ray.Dir.Dot(normal)
	if denom >= 0 {
		return -1
	}
	t = -(ray.Origin.Dot(normal) + dist) / denom
	return
}

func (p Plane) intersect(r g.Ray) float64 {
	return planeIntersect(p.Normal, p.Dist, r)
}

func (p Plane) traceTo(t float64, r g.Ray) (g.Vector, g.Vector) {
	return r.Incident(t), p.Normal
}

func (p Plane) sampleC(incident, normal, dir g.Vector, sc Scene) g.FloatColor {
	outColor := p.color
	_, modx := math.Modf(incident.X + 100)
	_, modz := math.Modf(incident.Z + 100)

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
		outColor = g.FloatColor{
			R: 128,
			G: 128,
			B: 128,
			A: 1.0,
		}
	} else {
		outColor = g.FloatColor{
			R: 255,
			G: 255,
			B: 255,
			A: 1.0,
		}
	}

	return outColor
}

func (p Plane) bounds(sceneBounds Aabb) Aabb {
	return sceneBounds
}
