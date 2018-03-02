package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	options := parseArgs(os.Args[1:])

	dx := vector{float64(2) / float64(options.width), 0, 0}
	dy := vector{
		0,
		float64(0.980580*-2) / float64(options.width),
		float64(0.196117*-2) / float64(options.height),
	}

	aperture := math.Sqrt(float64(options.width*options.height)) / 64.0

	samples := castQuad(
		vector{0, 2, -8},
		vector{-1, 1 + 0.980580, -3 + 0.196117},
		dx,
		dy,
		options.width,
		options.height,
		// []sampleMethod{createRgssSampler()},
		// []sampleMethod{},
		[]sampleMethod{createRgssSampler(), createDofSampler(5, 3, aperture)},
		// []sampleMethod{createDofSampler(5, 3, aperture)},
	)
	/*
		samples := castQuadDof(
			// vector{0, 1, -3},
			// vector{-1, 2, -2},
			vector{0, 0, -8},
			vector{-2, 2, -4},
			vector{float64(4) / float64(options.width), 0, 0},
			vector{0, float64(-4) / float64(options.height), 0},
			options.width,
			options.height,
			6, 3, 4,
		)
	*/

	scene := scene{
		shapes: []shape{
			sphere{
				material: defaultShader(color.RGBA{255, 0, 255, 255}, 0.8),
				center:   vector{0, .5, 2},
				radius:   .5,
			},
			sphere{
				createMaterial(color.RGBA{255, 255, 0, 255}, 0.3, 0.8, 0.7, 4),
				vector{-1.5, .4, 3},
				0.4,
			},
			sphere{
				createMaterial(color.RGBA{0, 255, 255, 255}, 0.9, 1.0, 1.0, 64),
				vector{1.5, .6, 4.5},
				0.6,
			},
			sphere{
				createMaterial(color.RGBA{128, 192, 255, 255}, 0.8, 1.0, 1.0, 32),
				vector{-1, .3, 5.5},
				0.3,
			},
			sphere{
				createMaterial(color.RGBA{255, 128, 192, 255}, 0.8, 1.0, 1.0, 32),
				vector{3, .7, 7.5},
				0.7,
			},
			sphere{
				createMaterial(color.RGBA{255, 192, 192, 255}, 0.8, 1.0, 1.0, 32),
				vector{.5, .3, 0},
				0.3,
			},
			sphere{
				createMaterial(color.RGBA{192, 255, 192, 255}, 0.8, 1.0, 1.0, 32),
				vector{-.5, .3, -1},
				0.3,
			},
			sphere{
				createMaterial(color.RGBA{192, 192, 255, 255}, 0.8, 1.0, 1.0, 32),
				vector{1, .3, -1.5},
				0.3,
			},
			plane{
				createMaterial(color.RGBA{192, 192, 192, 255}, 0.0, 1.0, 0.0, 0),
				vector{0, 1, 0},
				0,
			},
		},
		lights: []light{
			pointLight{vector{2, 2, -2}, 10},
			pointLight{vector{-1, 2, -3}, 5},
		},
		ambientLight: 0.3,
	}

	colors := make([]color.RGBA, options.width*options.height)
	// depths := make([]float64, options.width*options.height)
	batchSize := 128 * 128

	var wg sync.WaitGroup
	for i := 0; i < options.width*options.height; i += batchSize {
		wg.Add(1)
	}

	for i := 0; i < options.width*options.height; i += batchSize {
		go func(index int, endIndex int) {
			defer wg.Done()
			for i := index; i < endIndex && i < options.width*options.height; i++ {
				t := traceParams{
					10,
					1.0,
					color.RGBA{0, 0, 0, 0},
				}
				// colors[i], depths[i] = traceSample(samples[i], t, scene)
				colors[i], _ = traceSample(samples[i], t, scene)
			}
		}(i, i+batchSize)
	}

	wg.Wait()

	/*
		depthimage := image.NewRGBA(
			image.Rectangle{
				image.Point{0, 0},
				image.Point{options.width, options.height},
			},
		)
		for y := 0; y < options.height; y++ {
			for x := 0; x < options.width; x++ {
				depth := uint8(math.Min(255, math.Max(depths[y*options.width+x]*10, 0)))
				depthimage.Set(x, y, color.RGBA{
					depth,
					depth,
					depth,
					255,
				})
			}
		}

		depthfile, _ := os.Create("test_depths.png")
		png.Encode(depthfile, depthimage)
	*/

	testimage := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{options.width, options.height},
		},
	)

	for y := 0; y < options.height; y++ {
		for x := 0; x < options.width; x++ {
			testimage.Set(x, y, colors[y*options.width+x])
		}
	}

	testfile, _ := os.Create("test.png")
	png.Encode(testfile, testimage)

	elapsed := time.Since(start)
	fmt.Printf("Raytracing completed after %s\n", elapsed)
}
