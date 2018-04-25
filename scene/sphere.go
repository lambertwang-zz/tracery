package scene

import (
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

type sphere struct {
	Material
	center g.Vector
	radius float64
}

func (s sphere) shouldTest(r g.Ray) bool {
	return s.center.Minus(r.Origin).Dot(r.Dir) > 0
}

func (s sphere) intersect(r g.Ray) float64 {
	v := r.Origin.Minus(s.center)
	vd := v.Dot(r.Dir)
	discriminant := math.Pow(vd, 2) - (math.Pow(v.Magnitude(), 2) - math.Pow(s.radius, 2))

	// If the discriminant is negative, the ray does not intersect the sphere
	if discriminant < 0 {
		return -1
	}

	sqrtd := math.Sqrt(discriminant)
	// we are only concerned with the nearest point
	// return -vd + sqrtd, -vd - sqrtd
	return -vd - sqrtd
}

func (s sphere) traceTo(t float64, r g.Ray) (incident g.Vector, normal g.Vector) {
	// n = (y - c) / || y - c|| where y = s + td
	incident = r.Incident(t)
	normal = incident.Minus(s.center).Norm()
	return
}

func (s sphere) bounds(sceneBounds Aabb) Aabb {
	return Aabb{
		[3]Slab{
			Slab{-s.center.X - s.radius, -s.center.X + s.radius, g.Vector{X: 1, Y: 0, Z: 0}},
			Slab{-s.center.Y - s.radius, -s.center.Y + s.radius, g.Vector{X: 0, Y: 1, Z: 0}},
			Slab{-s.center.Z - s.radius, -s.center.Z + s.radius, g.Vector{X: 0, Y: 0, Z: 1}},
		},
		s.center,
	}
}
