package main

import (
	"math"
)

type vector struct {
	x float64
	y float64
	z float64
}

type ray struct {
	origin vector
	dir    vector
}

func addVectors(vecs ...vector) (out vector) {
	out = vector{0, 0, 0}
	for _, v := range vecs {
		out.x += v.x
		out.y += v.y
		out.z += v.z
	}
	return
}

func subtractVector(lhs vector, rhs vector) vector {
	return vector{lhs.x - rhs.x, lhs.y - rhs.y, lhs.z - rhs.z}
}

func dotProduct(lhs vector, rhs vector) float64 {
	return lhs.x*rhs.x + lhs.y*rhs.y + lhs.z*rhs.z
}

func (v vector) scale(scalar float64) vector {
	return vector{v.x * scalar, v.y * scalar, v.z * scalar}
}

func (v vector) magnitude() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2) + math.Pow(v.z, 2))
}

func (v vector) norm() vector {
	return v.scale(1.0 / v.magnitude())
}

func (v vector) neg() vector {
	return vector{-v.x, -v.y, -v.z}
}

func (r ray) incident(t float64) vector {
	return addVectors(r.origin, r.dir.scale(t))
}

func lineToRay(origin vector, target vector) ray {
	line := vector{target.x - origin.x, target.y - origin.y, target.z - origin.z}
	return ray{origin, line.norm()}
}
