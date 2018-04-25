package geometry

import "math"

// TMatrix3d TODO: Comment
//
// |	tf[0]	tf[1]	tf[2]	tl.X	|
// |	tf[3]	tf[4]	tf[5]	tl.Y	|
// |	tf[6]	tf[7]	tf[8]	tl.Z	|
// |	0		0		0		1		|
//
type TMatrix3d struct {
	tf [9]float64
	tl Vector
}

// Mult TODO: Why Why Why?
func (a TMatrix3d) Mult(b TMatrix3d) (c TMatrix3d) {
	c.tf[0], c.tf[1], c.tf[2],
		c.tf[3], c.tf[4], c.tf[5],
		c.tf[6], c.tf[7], c.tf[8],
		c.tl.X, c.tl.Y, c.tl.Z =
		a.tf[0]*b.tf[0]+a.tf[1]*b.tf[3]+a.tf[2]*b.tf[6],
		a.tf[0]*b.tf[1]+a.tf[1]*b.tf[4]+a.tf[2]*b.tf[7],
		a.tf[0]*b.tf[2]+a.tf[1]*b.tf[5]+a.tf[2]*b.tf[8],
		a.tf[3]*b.tf[0]+a.tf[4]*b.tf[3]+a.tf[5]*b.tf[6],
		a.tf[3]*b.tf[1]+a.tf[4]*b.tf[4]+a.tf[5]*b.tf[7],
		a.tf[3]*b.tf[2]+a.tf[4]*b.tf[5]+a.tf[5]*b.tf[8],
		a.tf[6]*b.tf[0]+a.tf[7]*b.tf[3]+a.tf[8]*b.tf[6],
		a.tf[6]*b.tf[1]+a.tf[7]*b.tf[4]+a.tf[8]*b.tf[7],
		a.tf[6]*b.tf[2]+a.tf[7]*b.tf[5]+a.tf[8]*b.tf[8],
		a.tf[0]*b.tl.X+a.tf[1]*b.tl.Y+a.tf[2]*b.tl.Z+a.tl.X,
		a.tf[3]*b.tl.X+a.tf[4]*b.tl.Y+a.tf[5]*b.tl.Z+a.tl.Y,
		a.tf[6]*b.tl.X+a.tf[7]*b.tl.Y+a.tf[8]*b.tl.Z+a.tl.Z

	return
}

// CreateIdentity TODO: Comment
func CreateIdentity() (c TMatrix3d) {
	c.tf[0], c.tf[4], c.tf[8] = 1, 1, 1
	return
}

// CreateTranslate TODO: Comment
func CreateTranslate(v Vector) (c TMatrix3d) {
	c.tl = v
	return
}

// CreateScale TODO: Comment
func CreateScale(f float64) (c TMatrix3d) {
	c.tf[0], c.tf[4], c.tf[8] = f, f, f
	return
}

// CreateRotate TODO: Comment
func CreateRotate(theta float64, v Vector) (c TMatrix3d) {
	cosTh, sinTh := math.Cos(theta), math.Sin(theta)
	iCosTh := 1 - cosTh
	lSinTh, mSinTh, nSinTh,
		lICosTh, mICosTh, nICosTh :=
		v.X*sinTh, v.Y*sinTh, v.Z*sinTh,
		v.X*iCosTh, v.Y*iCosTh, v.Z*iCosTh

	c.tf[0], c.tf[1], c.tf[2],
		c.tf[3], c.tf[4], c.tf[5],
		c.tf[6], c.tf[7], c.tf[8] =
		v.X*lICosTh+cosTh, v.Y*lICosTh-nSinTh, v.Z*lICosTh+mSinTh,
		v.X*mICosTh+nSinTh, v.Y*mICosTh+cosTh, v.Z*mICosTh-lSinTh,
		v.X*nICosTh-mSinTh, v.Y*nICosTh+lSinTh, v.Z*nICosTh+cosTh
	return
}

// Transform TODO: Comment
func (a TMatrix3d) Transform(b Vector) (c Vector) {
	c.X, c.Y, c.Z =
		a.tf[0]*b.X+a.tf[1]*b.Y+a.tf[2]*b.Z+a.tl.X,
		a.tf[3]*b.X+a.tf[4]*b.Y+a.tf[5]*b.Z+a.tl.Y,
		a.tf[6]*b.X+a.tf[7]*b.Y+a.tf[8]*b.Z+a.tl.Z
	return
}

// Translate TODO: Comment
func (a TMatrix3d) Translate(b Vector) (c TMatrix3d) {
	c = a.Mult(CreateTranslate(b))
	return
}

// Scale TODO: Comment
func (a TMatrix3d) Scale(b float64) (c TMatrix3d) {
	c = a.Mult(CreateScale(b))
	return
}

// Rotate TODO: Comment
func Rotate() {

}
