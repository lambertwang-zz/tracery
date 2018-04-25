package geometry

import (
	"math"
)

// Vector TODO: Comment
type Vector struct {
	X, Y, Z float64
}

// Ray TODO: Comment
type Ray struct {
	Origin, Dir Vector
}

// AddVectors TODO: Comment
func AddVectors(vecs ...Vector) (out Vector) {
	out = Vector{0, 0, 0}
	for _, v := range vecs {
		out = out.Plus(v)
	}
	return
}

func subtractVector(lhs Vector, rhs Vector) Vector {
	return lhs.Minus(rhs)
}

// DotProduct TODO: Comment
func DotProduct(lhs Vector, rhs Vector) float64 {
	return lhs.Dot(rhs)
}

// CrossProduct TODO: Comment
func CrossProduct(lhs Vector, rhs Vector) Vector {
	return Vector{
		lhs.Y*rhs.Z - lhs.Z*rhs.Y,
		lhs.Z*rhs.X - lhs.X*rhs.Z,
		lhs.X*rhs.Y - lhs.Y*rhs.X,
	}
}

// Dot TODO: Comment
func (v Vector) Dot(rhs Vector) float64 {
	return v.X*rhs.X + v.Y*rhs.Y + v.Z*rhs.Z
}

// Plus TODO: Comment
func (v Vector) Plus(rhs Vector) Vector {
	return Vector{v.X + rhs.X, v.Y + rhs.Y, v.Z + rhs.Z}
}

// Minus TODO: Comment
func (v Vector) Minus(rhs Vector) Vector {
	return Vector{v.X - rhs.X, v.Y - rhs.Y, v.Z - rhs.Z}
}

// Scale TODO: Comment
func (v Vector) Scale(scalar float64) Vector {
	return Vector{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

// Magnitude TODO: Comment
func (v Vector) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

// Norm TODO: Comment
func (v Vector) Norm() Vector {
	return v.Scale(1.0 / v.Magnitude())
}

// Neg TODO: Comment
func (v Vector) Neg() Vector {
	return Vector{-v.X, -v.Y, -v.Z}
}

// Incident TODO: Comment
func (r Ray) Incident(t float64) Vector {
	return r.Origin.Plus(r.Dir.Scale(t))
}

// LineToRay TODO: Comment
func LineToRay(origin Vector, target Vector) Ray {
	return Ray{origin, target.Minus(origin).Norm()}
}
