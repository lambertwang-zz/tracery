package tracer

import (
	g "github.com/lambertwang/tracery/geometry"
)

// Sample TODO: Comment
type Sample struct {
	dispatchedRays []g.Ray
}

// SampleMethod TODO: Comment
type SampleMethod func(
	g.Vector, g.Vector,
	g.Vector, g.Vector,
	*[]g.Ray)

// SampleSingle TODO: Comment
func SampleSingle(
	origin g.Vector, corner g.Vector,
	dx g.Vector, dy g.Vector,
	x int, y int,
	samplers []SampleMethod) (outSample Sample) {
	target := g.AddVectors(
		corner,
		dx.Scale(float64(x)),
		dy.Scale(float64(y)),
	)
	rays := []g.Ray{
		g.Ray{Origin: origin, Dir: target},
	}

	for _, sampler := range samplers {
		var newRays []g.Ray
		for _, r := range rays {
			sampler(r.Origin, r.Dir, dx, dy, &newRays)
		}
		rays = nil
		rays = make([]g.Ray, len(newRays))
		copy(rays, newRays)
	}

	for _, r := range rays {
		outSample.dispatchedRays = append(
			outSample.dispatchedRays,
			g.LineToRay(r.Origin, r.Dir),
		)
	}

	return
}

func sampleQuad(
	origin g.Vector, corner g.Vector,
	dx g.Vector, dy g.Vector,
	stepsX int, stepsY int,
	samplers []SampleMethod) []Sample {
	samples := make([]Sample, stepsX*stepsY)

	for y := 0; y < stepsY; y++ {
		for x := 0; x < stepsX; x++ {
			samples[y*stepsX+x] = SampleSingle(origin, corner, dx, dy, x, y, samplers)
		}
	}

	return samples
}
