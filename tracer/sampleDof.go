package tracer

import (
	"math/rand"

	g "github.com/lambertwang/tracery/geometry"
)

// CreateDofSampler TODO: Comment
func CreateDofSampler(focalDist float64, samples int, aperture float64) SampleMethod {
	return func(
		origin g.Vector, target g.Vector,
		dx g.Vector, dy g.Vector,
		rays *[]g.Ray) {
		var imagePlanePoints []g.Vector
		for i := 0; i < samples; i++ {
			for j := 0; j < samples; j++ {
				imagePlanePoints = append(imagePlanePoints,
					g.AddVectors(
						target,
						dx.Scale(
							aperture*
								(float64(i)-float64(samples)/2.0+rand.Float64())),
						dy.Scale(
							aperture*
								(float64(j)-float64(samples)/2.0+rand.Float64())),
					),
				)
			}
		}

		focalPoint := g.AddVectors(target.Minus(origin).Norm().Scale(focalDist), target)
		for _, p := range imagePlanePoints {
			*rays = append(*rays, g.Ray{Origin: p, Dir: focalPoint})
		}
	}
}
