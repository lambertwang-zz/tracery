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
	dy := vector{
		0,
		float64(0.980580*-2) / float64(options.width),
		float64(0.196117*-2) / float64(options.width),
	}

	rayCasts, samples := castQuadRgss(
		vector{0, 3.5, -8},
		vector{-1, 1.5 + 0.980580, -3 + 0.196117},
		dx,
		dy,
		options.width,
		options.height,
	)

	scene := scene{
		[]shape{
			sphere{
				material{
					color.RGBA{255, 0, 255, 255},
					0.7,
				},
				vector{.25, .6, 0},
				0.6,
			},
			sphere{
				material{
					color.RGBA{255, 255, 0, 255},
					0.3,
				},
				vector{-.5, .5, -1},
				0.5,
			},
			sphere{
				material{
					color.RGBA{0, 255, 255, 255},
					0.9,
				},
				vector{1, .4, -.75},
				0.4,
			},
			sphere{
				material{
					color.RGBA{255, 192, 192, 255},
					0.8,
				},
				vector{.5, .3, -1.5},
				0.3,
			},
			plane{
				material{
					color.RGBA{64, 64, 64, 64},
					0,
				},
				vector{0, 1, 0},
				0,
			},
		},
		[]light{
			pointLight{vector{2, 2, -2}},
			pointLight{vector{-1, 2, -3}},
		},
	}

	colors := make([]color.RGBA, options.width*options.height)
	sampleColors := make([]color.RGBA, options.width*options.height*samples)
	batchSize := 128

	var wg sync.WaitGroup
	for y := 0; y < options.height; y += batchSize {
		for x := 0; x < options.width; x += batchSize {
			wg.Add(1)
		}
	}

	for y := 0; y < options.height; y += batchSize {
		for x := 0; x < options.width; x += batchSize {
			go func(initX int, initY int, endX int, endY int) {
				defer wg.Done()
				for y := initY; y < endY && y < options.height; y++ {
					for x := initX; x < endX && x < options.width; x++ {
						t := traceParams{
							10,
							1.0,
							color.RGBA{0, 0, 0, 0},
						}
						c := make([]color.RGBA, samples)
						c[0] = trace(rayCasts[(y*options.width+x)*samples], t, scene)
						c[1] = trace(rayCasts[(y*options.width+x)*samples+1], t, scene)
						c[2] = trace(rayCasts[(y*options.width+x)*samples+2], t, scene)
						c[3] = trace(rayCasts[(y*options.width+x)*samples+3], t, scene)
						sampleColors[(y*options.width+x)*samples] = c[0]
						sampleColors[(y*options.width+x)*samples+1] = c[1]
						sampleColors[(y*options.width+x)*samples+2] = c[2]
						sampleColors[(y*options.width+x)*samples+3] = c[3]
						colors[y*options.width+x] = meanColor(c[0], c[1], c[2], c[3])
					}
				}
			}(x, y, x+batchSize, y+batchSize)
		}
	}

	wg.Wait()

	imageSample := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{options.width * 2, options.height * 2},
		},
	)
	for y := 0; y < options.height; y++ {
		for x := 0; x < options.width; x++ {
			imageSample.Set(x, y, sampleColors[(y*options.width+x)*samples])
			imageSample.Set(x+1, y, sampleColors[(y*options.width+x)*samples+1])
			imageSample.Set(x, y+1, sampleColors[(y*options.width+x)*samples+2])
			imageSample.Set(x+1, y+1, sampleColors[(y*options.width+x)*samples+3])
		}
	}

	image := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{options.width, options.height},
		},
	)
	for y := 0; y < options.height; y++ {
		for x := 0; x < options.width; x++ {
			image.Set(x, y, colors[y*options.width+x])
		}
	}

	file, _ := os.Create("test.png")
	png.Encode(file, image)

	file2, _ := os.Create("testSamples.png")
	png.Encode(file2, imageSample)

	elapsed := time.Since(start)
	fmt.Printf("Raytracing completed after %s\n", elapsed)
}
