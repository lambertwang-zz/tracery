package main

// Triangle handedness is CCW
type triangle struct {
	material
	p0, p1, p2 vector
	// Plane properties
	n0 vector
	n1 vector
	d0 float64
	d1 float64
}

func (t triangle) shouldTest(ray ray) bool {
	return true
}

func createTriangle(m material, p0 vector, p1 vector, p2 vector) triangle {
	t := triangle{
		m,
		p0, p1, p2,
		vector{0, 0, 0},
		vector{0, 0, 0},
		0, 0,
	}
	t.n0, t.n1 = t.norm()
	t.d0 = -dotProduct(t.n0, t.p0)
	t.d1 = -dotProduct(t.n1, t.p0)
	return t
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
	e1 := subtractVector(incident, t.p1)
	e2 := subtractVector(incident, t.p2)

	if dotProduct(normal, crossProduct(a, e0)) > 0 &&
		dotProduct(normal, crossProduct(b, e1)) > 0 &&
		dotProduct(normal, crossProduct(c, e2)) > 0 {
		return planeIntersection
	}

	return -1
}

func (t triangle) reflect(x ray, t0 float64, sc scene) (vector, []vector, floatColor, float64, vector) {
	outColor := t.color
	planeIntersection := planeIntersect(t.n0, t.d0, x)
	normal := t.n0
	if planeIntersection < 0 {
		normal = t.n1
	}
	incident := x.incident(t0)
	reflection := subtractVector(x.dir, normal.scale(2*(dotProduct(normal, x.dir))))

	return incident, []vector{reflection}, outColor, t.reflectance, normal
}

func (t triangle) getMaterial() material {
	return t.material
}
