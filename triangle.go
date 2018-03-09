package main

import "math"

// Triangle handedness is CCW
type triangle struct {
	material
	p0, p1, p2 vector
	// Plane properties
	n0 vector
	n1 vector
	d0 float64
	d1 float64

	vertexNormals [3]vector
}

func createTriangle(m material, p0 vector, p1 vector, p2 vector) triangle {
	t := triangle{
		m,
		p0, p1, p2,
		vector{0, 0, 0},
		vector{0, 0, 0},
		0, 0,
		[3]vector{
			vector{0, 0, 0},
			vector{0, 0, 0},
			vector{0, 0, 0},
		},
	}
	t.n0, t.n1 = t.norm()
	t.vertexNormals = [3]vector{t.n0, t.n0, t.n0}
	t.d0 = -dotProduct(t.n0, t.p0)
	t.d1 = -dotProduct(t.n1, t.p0)
	return t
}

func (t triangle) computeBarycentric(p vector) (coords vector) {
	v0 := t.p1.minus(t.p0)
	v1 := t.p2.minus(t.p0)
	v2 := p.minus(t.p0)

	d00 := dotProduct(v0, v0)
	d01 := dotProduct(v0, v1)
	d11 := dotProduct(v1, v1)
	d20 := dotProduct(v2, v0)
	d21 := dotProduct(v2, v1)

	denom := d00*d11 - d01*d01

	v := (d11*d20 - d01*d21) / denom
	w := (d00*d21 - d01*d20) / denom
	u := 1.0 - v - w
	coords = vector{u, v, w}

	return
}

func (t triangle) computeSurfaceProps(p vector) (normal vector) {
	uvw := t.computeBarycentric(p)
	normal = t.vertexNormals[0].scale(uvw.x).plus(
		t.vertexNormals[1].scale(uvw.y).plus(
			t.vertexNormals[2].scale(uvw.z)))
	return
}

func (t triangle) norm() (vector, vector) {
	a := subtractVector(t.p1, t.p0)
	b := subtractVector(t.p2, t.p0)
	return crossProduct(a, b).norm(), crossProduct(b, a).norm()
}

func (t triangle) intersect(ray ray) float64 {
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

	incident := ray.incident(planeIntersection)

	var a, b, c vector
	if useBackface {
		a = subtractVector(t.p2, t.p0)
		b = subtractVector(t.p0, t.p1)
		c = subtractVector(t.p1, t.p2)
	} else {
		a = subtractVector(t.p1, t.p0)
		b = subtractVector(t.p2, t.p1)
		c = subtractVector(t.p0, t.p2)
	}

	e0 := subtractVector(incident, t.p0)
	if dotProduct(normal, crossProduct(a, e0)) < 0 {
		return -1
	}
	e1 := subtractVector(incident, t.p1)
	if dotProduct(normal, crossProduct(b, e1)) < 0 {
		return -1
	}
	e2 := subtractVector(incident, t.p2)
	if dotProduct(normal, crossProduct(c, e2)) < 0 {
		return -1
	}

	return planeIntersection
}

func (t triangle) traceTo(t0 float64, params traceParams) (incident vector, normal vector) {
	incident = params.rayCast.incident(t0)
	// planeIntersection := planeIntersect(t.n0, t.d0, params.rayCast)
	normal = t.computeSurfaceProps(incident)
	// normal = t.n0
	// if planeIntersection < 0 {
	// 	normal = t.n1
	// }

	return
}

func (t triangle) bounds(sceneBounds aabb) aabb {
	return aabb{
		[3]slab{
			slab{
				-math.Max(t.p0.x, math.Max(t.p1.x, t.p2.x)),
				-math.Min(t.p0.x, math.Min(t.p1.x, t.p2.x)),
				vector{1, 0, 0},
			},
			slab{
				-math.Max(t.p0.y, math.Max(t.p1.y, t.p2.y)),
				-math.Min(t.p0.y, math.Min(t.p1.y, t.p2.y)),
				vector{0, 1, 0},
			},
			slab{
				-math.Max(t.p0.z, math.Max(t.p1.z, t.p2.z)),
				-math.Min(t.p0.z, math.Min(t.p1.z, t.p2.z)),
				vector{0, 0, 1},
			},
		},
		t.p0.plus(t.p1.plus(t.p2)).scale(1 / 3),
	}
}
