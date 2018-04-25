package scene

import (
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// Triangle handedness is CCW
type Triangle struct {
	Material
	p0, p1, p2 g.Vector
	// Plane properties
	n0, n1 g.Vector
	d0, d1 float64

	VertexNormals [3]g.Vector
}

// CreateTriangle TODO: Comment
func CreateTriangle(m Material, p0, p1, p2 g.Vector) (t Triangle) {
	t.Material = m
	t.p0 = p0
	t.p1 = p1
	t.p2 = p2
	t.n0, t.n1 = t.norm()
	t.VertexNormals = [3]g.Vector{t.n0, t.n0, t.n0}
	t.d0 = -t.n0.Dot(t.p0)
	t.d1 = -t.n1.Dot(t.p0)
	return t
}

func (t Triangle) computeBarycentric(p g.Vector) (coords g.Vector) {
	v0 := t.p1.Minus(t.p0)
	v1 := t.p2.Minus(t.p0)
	v2 := p.Minus(t.p0)

	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)

	denom := d00*d11 - d01*d01

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom
	u := 1.0 - v - w
	coords = g.Vector{X: u, Y: v, Z: w}

	return
}

func (t Triangle) computeSurfaceProps(p g.Vector) (normal g.Vector) {
	uvw := t.computeBarycentric(p)
	normal = t.VertexNormals[0].Scale(uvw.X).Plus(
		t.VertexNormals[1].Scale(uvw.Y).Plus(
			t.VertexNormals[2].Scale(uvw.Z))).Norm()
	return
}

func (t Triangle) norm() (g.Vector, g.Vector) {
	a := t.p1.Minus(t.p0)
	b := t.p2.Minus(t.p0)
	return g.CrossProduct(a, b).Norm(), g.CrossProduct(b, a).Norm()
}

func (t Triangle) intersect(ray g.Ray) float64 {
	planeIntersection := planeIntersect(t.n0, t.d0, ray)
	normal := t.n0
	useBackface := false
	if planeIntersection < 0 {
		normal = t.n1
		planeIntersection = planeIntersect(t.n1, t.d1, ray)
		useBackface = true
	}
	if planeIntersection < 0 {
		return -1
	}

	incident := ray.Incident(planeIntersection)

	var a, b, c g.Vector
	if useBackface {
		a = t.p2.Minus(t.p0)
		b = t.p0.Minus(t.p1)
		c = t.p1.Minus(t.p2)
	} else {
		a = t.p1.Minus(t.p0)
		b = t.p2.Minus(t.p1)
		c = t.p0.Minus(t.p2)
	}

	e0 := incident.Minus(t.p0)
	if normal.Dot(g.CrossProduct(a, e0)) < 0 {
		return -1
	}
	e1 := incident.Minus(t.p1)
	if normal.Dot(g.CrossProduct(b, e1)) < 0 {
		return -1
	}
	e2 := incident.Minus(t.p2)
	if normal.Dot(g.CrossProduct(c, e2)) < 0 {
		return -1
	}

	return planeIntersection

}

func (t Triangle) traceTo(t0 float64, r g.Ray) (incident g.Vector, normal g.Vector) {
	incident = r.Incident(t0)
	// normal = t.computeSurfaceProps(incident)

	planeIntersection := planeIntersect(t.n0, t.d0, r)
	normal = t.n0
	if planeIntersection < 0 {
		// perpendicular := g.Vector{X: 0, Y: -t.n0.Z, Z: -t.n0.Y}.Norm()
		// normal = normal.Minus(perpendicular.Scale(2 * perpendicular.Dot(normal)))
		normal = t.n1
	}

	return
}

func (t Triangle) bounds(sceneBounds Aabb) Aabb {
	return Aabb{
		[3]Slab{
			Slab{
				-math.Max(t.p0.X, math.Max(t.p1.X, t.p2.X)),
				-math.Min(t.p0.X, math.Min(t.p1.X, t.p2.X)),
				g.Vector{X: 1, Y: 0, Z: 0},
			},
			Slab{
				-math.Max(t.p0.Y, math.Max(t.p1.Y, t.p2.Y)),
				-math.Min(t.p0.Y, math.Min(t.p1.Y, t.p2.Y)),
				g.Vector{X: 0, Y: 1, Z: 0},
			},
			Slab{
				-math.Max(t.p0.Z, math.Max(t.p1.Z, t.p2.Z)),
				-math.Min(t.p0.Z, math.Min(t.p1.Z, t.p2.Z)),
				g.Vector{X: 0, Y: 0, Z: 1},
			},
		},
		t.p0.Plus(t.p1.Plus(t.p2)).Scale(1.0 / 3.0),
	}
}
