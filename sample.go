package main

import "math/rand"

type sample struct {
	dispatchedRays []ray
}

type sampleMethod func(vector, vector, vector, vector, *[]ray)

func sampleSingle(origin vector, corner vector, dx vector, dy vector, x int, y int, samplers []sampleMethod) (outSample sample) {
	target := addVectors(
		corner,
		dx.scale(float64(x)),
		dy.scale(float64(y)),
	)
	rays := []ray{
		ray{origin, target},
	}

	for _, sampler := range samplers {
		var newRays []ray
		for _, r := range rays {
			sampler(r.origin, r.dir, dx, dy, &newRays)
		}
		rays = nil
		rays = make([]ray, len(newRays))
		copy(rays, newRays)
	}

	for _, r := range rays {
		outSample.dispatchedRays = append(
			outSample.dispatchedRays,
			lineToRay(r.origin, r.dir),
		)
	}

	return
}

func sampleQuad(origin vector, corner vector, dx vector, dy vector, stepsX int, stepsY int, samplers []sampleMethod) []sample {
	samples := make([]sample, stepsX*stepsY)

	for y := 0; y < stepsY; y++ {
		for x := 0; x < stepsX; x++ {
			samples[y*stepsX+x] = sampleSingle(origin, corner, dx, dy, x, y, samplers)
		}
	}

	return samples
}

func createRgssSampler() sampleMethod {
	return func(origin vector, target vector, dx vector, dy vector, rays *[]ray) {
		*rays = append(
			*rays,
			ray{
				origin,
				addVectors(target,
					dx.scale(.125),
					dy.scale(.375),
				),
			},
			ray{
				origin,
				addVectors(target,
					dx.scale(.625),
					dy.scale(.125),
				),
			},
			ray{
				origin,
				addVectors(target,
					dx.scale(.375),
					dy.scale(.875),
				),
			},
			ray{
				origin,
				addVectors(target,
					dx.scale(.875),
					dy.scale(.625),
				),
			},
		)
	}
}

func createDofSampler(focalDist float64, samples int, aperture float64) sampleMethod {
	return func(origin vector, target vector, dx vector, dy vector, rays *[]ray) {
		var imagePlanePoints []vector
		for i := 0; i < samples; i++ {
			for j := 0; j < samples; j++ {
				imagePlanePoints = append(imagePlanePoints,
					addVectors(
						target,
						dx.scale(
							aperture*
								(float64(i)-float64(samples)/2.0+rand.Float64())),
						dy.scale(
							aperture*
								(float64(j)-float64(samples)/2.0+rand.Float64())),
					),
				)
			}
		}

		focalPoint := addVectors(subtractVector(target, origin).norm().scale(focalDist), target)
		for _, p := range imagePlanePoints {
			*rays = append(*rays, ray{p, focalPoint})
		}
	}
}
