package scene

import (
	"image/color"
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// Material properties for rendering
type Material struct {
	color       g.FloatColor
	reflectance float64

	// Shading properties
	// alpha    float64
	diffuse  float64
	specular float64
	hardness float64

	refractionIndex float64
}

// DefaultMaterial TODO: Comment
func DefaultMaterial() Material {
	return Material{
		g.FloatColor{R: 192, G: 192, B: 192, A: 1.0},
		0.0, 1.0, 1.0, 32, 0.0,
	}
}

// DefaultShader TODO: Comment
func DefaultShader(c color.RGBA, r float64) Material {
	return Material{g.RgbaToFloat(c), r, 1.0, 1.0, 32, 0.0}
}

// CreateMaterial TODO: Comment
func CreateMaterial(c color.RGBA, r float64, diff float64, spec float64, hard float64, refr float64) Material {
	return Material{g.RgbaToFloat(c), r, diff, spec, hard, refr}
}

func (m Material) getMaterial() Material {
	return m
}

func (m Material) shouldTest(g.Ray) bool {
	return true
}

func (m Material) sampleC(incident, normal, dir g.Vector, sc Scene) g.FloatColor {
	return m.color
}

func (m Material) reflect(incident, normal, dir g.Vector, sc Scene) (reflection g.Vector) {
	reflection = dir.Minus(normal.Scale(2 * normal.Dot(dir)))
	return
}

func (m Material) refract(incident, normal, dir g.Vector, sc Scene) (refraction g.Vector) {
	refractionIndex := 1.0 / m.refractionIndex
	/*
		theta1 := math.Acos(dotProduct(dir, normal))
		theta2 := math.Asin(math.Sin(theta1) * refractionIndex)
		refraction := subtractVector(
			addVectors(dir.scale(math.Sin(theta2)/math.Sin(theta1)), normal.scale(math.Cos(theta1)*math.Sin(theta2)/math.Sin(theta1))),
			normal.scale(math.Cos(theta2)),
		).norm()
	*/

	c1 := normal.Dot(dir)
	c2 := math.Sqrt(1 - math.Pow(refractionIndex, 2)*(1-math.Pow(c1, 2)))
	refraction = dir.Scale(refractionIndex).Plus(normal.Scale(refractionIndex*c1 - c2)).Norm()
	return
}
