package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var testImage *image.RGBA
var opts options

type mainWindow struct {
	*walk.MainWindow
	paintWidget *walk.CustomWidget
}

func (mw *mainWindow) renderImage(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	bmp, _ := walk.NewBitmapFromImage(testImage)
	canvas.DrawBitmapPart(
		bmp,
		walk.Rectangle{X: 0, Y: 0, Width: opts.width, Height: opts.height},
		walk.Rectangle{X: 0, Y: 0, Width: opts.width, Height: opts.height},
	)
	return nil
}

func main() {
	start := time.Now()
	opts = parseArgs(os.Args[1:])

	testImage = image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{opts.width, opts.height},
		},
	)
	mw := new(mainWindow)

	scene := scene{
		shapes: []shape{
			sphere{
				// material: createMaterial(color.RGBA{255, 0, 255, 64}, 0.4, 1.0, 1.0, 32, 1.5),
				material: createMaterial(color.RGBA{0, 0, 0, 0}, 0, 1.0, 1.0, 32, 1.5),
				center:   vector{0, .5, 2},
				radius:   .5,
			},
			sphere{
				createMaterial(color.RGBA{255, 255, 0, 255}, 0.3, 0.8, 0.7, 4, 0.0),
				vector{-1, .4, 3},
				0.4,
			},
			sphere{
				createMaterial(color.RGBA{0, 255, 255, 255}, 0.7, 1.0, 1.0, 64, 0.0),
				vector{1, .6, 4.5},
				0.6,
			},
			sphere{
				createMaterial(color.RGBA{128, 192, 255, 255}, 0.8, 1.0, 1.0, 32, 0.0),
				vector{-.5, .3, 5.5},
				0.3,
			},
			sphere{
				createMaterial(color.RGBA{255, 128, 192, 255}, 0.8, 1.0, 1.0, 32, 0.0),
				vector{2, .7, 7.5},
				0.7,
			},
			sphere{
				createMaterial(color.RGBA{255, 192, 192, 255}, 0.8, 1.0, 1.0, 32, 0.0),
				vector{.2, .3, 0},
				0.3,
			},
			sphere{
				createMaterial(color.RGBA{192, 255, 192, 255}, 0.8, 1.0, 1.0, 32, 0.0),
				vector{-.3, .3, -1},
				0.3,
			},
			sphere{
				createMaterial(color.RGBA{192, 192, 255, 255}, 0.8, 1.0, 1.0, 32, 0.0),
				vector{0.6, .3, -1.5},
				0.3,
			},
			plane{
				createMaterial(color.RGBA{192, 192, 192, 255}, 0.3, 1.0, 0.0, 0, 0.0),
				vector{0, 1, 0},
				0,
			},
		},
		lights: []light{
			pointLight{
				vector{4, 5, 4}, 10,
				floatColor{255, 192, 0, 1.0},
			},
			pointLight{
				vector{-3, 3, -5}, 8,
				floatColor{0, 255, 192, 1.0},
			},
		},
		ambientLight: 0.1,
	}
	/*
		dx := vector{float64(2) / float64(options.width), 0, 0}
		dy := vector{0, 0, float64(-2) / float64(options.width)}
		// aperture := math.Sqrt(float64(options.width*options.height)) / 64.0

		samples := castQuad(
			vector{0, 20, 0},
			vector{-1, 17, 1},
			dx,
			dy,
			options.width,
			options.height,
			[]sampleMethod{},
			// []sampleMethod{createRgssSampler()},
			// []sampleMethod{createDofSampler(5, 3, aperture)},
			// []sampleMethod{createRgssSampler(), createDofSampler(5, 3, aperture)},
		)

		scene := scene{
			shapes: []shape{
				sphere{
					createMaterial(color.RGBA{128, 128, 128, 0}, 0.0, 0.0, 1.0, 64, 1.3),
					vector{0, 3, 0},
					1,
				},
				plane{
					createMaterial(color.RGBA{192, 192, 192, 255}, 0.0, 1.0, 1.0, 64, 1.0),
					vector{0, 1, 0},
					0,
				},
			},
			lights: []light{
				pointLight{vector{5, 20, 0}, 18},
				pointLight{vector{-3, 20, 4}, 18},
				pointLight{vector{-4, 20, -3}, 18},
			},
			ambientLight: 0.3,
		}
	*/

	// aperture := math.Sqrt(float64(options.width*options.height)) / 128.0

	castOrigin := vector{0, 2, -8}
	castCorner := vector{-1, 1 + 0.980580, -3 + 0.196117}
	dx := vector{float64(2) / float64(opts.width), 0, 0}
	dy := vector{
		0,
		float64(0.980580*-2) / float64(opts.width),
		float64(0.196117*-2) / float64(opts.height),
	}
	samplers := []sampleMethod{}
	// []sampleMethod{createDofSampler(5, 2, aperture)},
	// samplers = []sampleMethod{createRgssSampler(), createDofSampler(5, 3, 8)}
	samplers = []sampleMethod{createRgssSampler()}
	colors := make([]color.RGBA, opts.width*opts.height)

	var wg sync.WaitGroup

	batchSize := 128
	for j := 0; j < opts.height; j += batchSize {
		for i := 0; i < opts.width; i += batchSize {
			wg.Add(1)
		}
	}

	for j := 0; j < opts.height; j += batchSize {
		for i := 0; i < opts.width; i += batchSize {
			go func(startX int, startY int, batch int) {
				defer wg.Done()
				for y := startY; y < startY+batch && y < opts.height; y++ {
					for x := startX; x < startX+batch && x < opts.width; x++ {
						t := traceParams{
							ray{vector{0, 0, 0}, vector{0, 0, 0}},
							5, 1.0,
						}
						targetSample := sampleSingle(castOrigin, castCorner, dx, dy, x, y, samplers)
						colors[y*opts.width+x], _ = traceSample(targetSample, t, scene)
						testImage.Set(x, y, colors[y*opts.width+x])
					}
					// Update the widget
					if mw.paintWidget != nil {
						mw.paintWidget.Synchronize(func() {
							mw.paintWidget.SetPaintMode(walk.PaintNoErase)
							mw.paintWidget.Invalidate()
						})
					}
				}
			}(i, j, batchSize)
		}
	}

	go func() {
		// Wait for all rendering batches to finish
		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("Raytracing completed after %s\n", elapsed)
	}()

	MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "Tracery preview",
		MinSize: Size{
			Height: opts.width,
			Width:  opts.height,
		},
		Size: Size{
			Height: opts.width,
			Width:  opts.height,
		},
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			CustomWidget{
				AssignTo:            &mw.paintWidget,
				ClearsBackground:    true,
				InvalidatesOnResize: true,
				Paint:               mw.renderImage,
			},
		},
	}.Run()

	wg.Wait()

	testfile, _ := os.Create("test.png")
	png.Encode(testfile, testImage)
}
