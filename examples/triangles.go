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

	overlap := 3.0
	shapes := []shape{
		sphere{
			createMaterial(color.RGBA{255, 128, 192, 255}, 0.8, 1.0, 1.0, 64, 1.0),
			Vector{0, 0, 8},
			.5,
		},
		createTriangle(
			createMaterial(color.RGBA{192, 192, 255, 255}, 0, 1.0, 1.0, 32, 0),
			Vector{-2, 0, 8 + overlap},
			Vector{-2, 2, 8},
			Vector{1, 1, 8 - overlap},
		),
		createTriangle(
			createMaterial(color.RGBA{255, 192, 255, 255}, 0, 1.0, 1.0, 32, 0),
			Vector{0, 2, 8 + overlap},
			Vector{2, 2, 8},
			Vector{1, -1, 8 - overlap},
		),
		createTriangle(
			createMaterial(color.RGBA{255, 255, 192, 255}, 0, 1.0, 1.0, 32, 0),
			Vector{2, 0, 8 + overlap},
			Vector{2, -2, 8},
			Vector{-1, -1, 8 - overlap},
		),
		createTriangle(
			createMaterial(color.RGBA{192, 255, 255, 255}, 0, 1.0, 1.0, 32, 0),
			Vector{0, -2, 8 + overlap},
			Vector{-2, -2, 8},
			Vector{-1, 1, 8 - overlap},
		),
	}

	// cubeModel := loadObjModel("cube.obj")
	/*
		cubeModel := loadObjModel("teapot.obj")
		for _, tri := range cubeModel.toShapes() {
			shapes = append(shapes, tri)
		}
	*/

	sceneBounds := Aabb{
		[3]Slab{
			Slab{-1000, 1000, Vector{1, 0, 0}},
			Slab{-1000, 1000, Vector{0, 1, 0}},
			Slab{-1000, 1000, Vector{0, 0, 1}},
		},
		Vector{0, 0, 0},
	}

	scene := scene{
		shapes: shapes,
		lights: []light{
			pointLight{
				Vector{2, 2, 0}, 10,
				floatColor{255, 192, 0, 1.0},
			},
			pointLight{
				Vector{-3, 1, 0}, 8,
				floatColor{0, 255, 192, 1.0},
			},
		},
		ambientLight: 0.3,
		bvh:          constructHeirarchy(&shapes, sceneBounds),
	}

	castOrigin := Vector{0, 0, 0}
	castCorner := Vector{-1, 1, 3}
	dx := Vector{float64(2) / float64(opts.width), 0, 0}
	dy := Vector{0, float64(-2) / float64(opts.width), 0}
	// aperture := math.Sqrt(float64(options.width*options.height)) / 64.0

	samplers := []sampleMethod{}
	// []sampleMethod{createDofSampler(5, 2, aperture)},
	// samplers = []sampleMethod{createRgssSampler(), createDofSampler(5, 3, 8)}
	// samplers = []sampleMethod{createRgssSampler()}

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
							ray{Vector{0, 0, 0}, Vector{0, 0, 0}},
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
