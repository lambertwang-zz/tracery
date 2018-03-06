package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	options := parseArgs(os.Args[1:])

	dx := vector{float64(2) / float64(options.width), 0, 0}
	dy := vector{0, float64(-2) / float64(options.width), 0}
	// aperture := math.Sqrt(float64(options.width*options.height)) / 64.0

	samples := castQuad(
		vector{0, 0, 0},
		vector{-1, 1, 3},
		dx,
		dy,
		options.width,
		options.height,
		[]sampleMethod{},
		// []sampleMethod{createRgssSampler()},
		// []sampleMethod{createDofSampler(5, 3, aperture)},
		// []sampleMethod{createRgssSampler(), createDofSampler(5, 3, aperture)},
	)

	var overlap float64 = 3
	scene := scene{
		shapes: []shape{
			/*
				createTriangle(
					createMaterial(color.RGBA{255, 255, 192, 255}, 0, 1.0, 1.0, 32),
					vector{-2, 0, 8 + overlap},
					vector{2, 2, 8 - overlap},
					vector{2, -2, 8 - overlap},
				),
				createTriangle(
					createMaterial(color.RGBA{192, 255, 255, 255}, 0, 1.0, 1.0, 32),
					vector{-2, 1, 9},
					vector{-2, -1, 9},
					vector{-2, 0, 7},
				),
			*/
			sphere{
				createMaterial(color.RGBA{255, 128, 192, 255}, 0.8, 0.0, 1.0, 64),
				vector{0, 0, 8},
				.5,
			},
			createTriangle(
				createMaterial(color.RGBA{192, 192, 255, 255}, 0, 1.0, 1.0, 32),
				vector{-2, 0, 8 + overlap},
				vector{-2, 2, 8},
				vector{1, 1, 8 - overlap},
			),
			createTriangle(
				createMaterial(color.RGBA{255, 192, 255, 255}, 0, 1.0, 1.0, 32),
				vector{0, 2, 8 + overlap},
				vector{2, 2, 8},
				vector{1, -1, 8 - overlap},
			),
			createTriangle(
				createMaterial(color.RGBA{255, 255, 192, 255}, 0, 1.0, 1.0, 32),
				vector{2, 0, 8 + overlap},
				vector{2, -2, 8},
				vector{-1, -1, 8 - overlap},
			),
			createTriangle(
				createMaterial(color.RGBA{192, 255, 255, 255}, 0, 1.0, 1.0, 32),
				vector{0, -2, 8 + overlap},
				vector{-2, -2, 8},
				vector{-1, 1, 8 - overlap},
			),
		},
		lights: []light{
			pointLight{vector{2, 2, 0}, 10},
			pointLight{vector{-3, 1, 0}, 8},
		},
		ambientLight: 0.3,
	}

	colors := make([]color.RGBA, options.width*options.height)
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
				colors[i], _ = traceSample(samples[i], t, scene)
			}
		}(i, i+batchSize)
	}

	wg.Wait()

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
