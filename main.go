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

	g "github.com/lambertwang/tracery/geometry"
	sc "github.com/lambertwang/tracery/scene"
	tr "github.com/lambertwang/tracery/tracer"

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
	defer bmp.Dispose()
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

	shapes := []sc.Shape{
	/*
		sc.Plane{
			sc.CreateMaterial(color.RGBA{192, 192, 192, 255}, 0.3, 1.0, 0.0, 0, 0.0),
			g.Vector{Y: 1},
			0,
		},
	*/
	}

	// cubeModel := loadObjModel("models/teapot.obj")
	cubeModel := loadObjModel("models/bunny.obj")
	for _, tri := range cubeModel.toShapes(
		sc.CreateMaterial(color.RGBA{255, 128, 192, 255}, .6, 1.0, 1.0, 32, 0.0),
		g.CreateRotate(math.Pi, g.Vector{X: 0, Y: 1, Z: 0}).Scale(2),
		// g.CreateScale(0.05),
		// g.CreateIdentity(),
	) {
		// for _, tri := range cubeModel.transformShapes(g.CreateRotate(math.Pi, g.Vector{X: 0, Y: 1, Z: 0})) {
		shapes = append(shapes, tri)
	}

	sceneBounds := sc.Aabb{
		Slabs: [3]sc.Slab{
			sc.Slab{-1000, 1000, g.Vector{1, 0, 0}},
			sc.Slab{-1000, 1000, g.Vector{0, 1, 0}},
			sc.Slab{-1000, 1000, g.Vector{0, 0, 1}},
		},
		Mean: g.Vector{0, 0, 0},
	}

	lights := []sc.Light{
		sc.PointLight{
			g.Vector{X: 16, Y: 10, Z: -10}, 30,
			g.FloatColor{R: 255, G: 192, B: 0, A: 1.0},
		},
		sc.PointLight{
			g.Vector{-15, 11, -10}, 28,
			g.FloatColor{0, 255, 192, 1.0},
		},
	}

	scene := sc.CreateScene(&shapes, &lights, 0.3, sceneBounds)

	castOrigin := g.Vector{0.001, 1, -7}
	castCorner := g.Vector{-1.001, 2, -5}
	dx := g.Vector{float64(2) / float64(opts.width), 0, 0}
	dy := g.Vector{0, float64(-2) / float64(opts.width), 0}
	// aperture := math.Sqrt(float64(options.width*options.height)) / 64.0

	samplers := []tr.SampleMethod{}
	// []sampleMethod{createDofSampler(5, 2, aperture)},
	// samplers = []sampleMethod{createRgssSampler(), createDofSampler(9, 3, 16)}
	samplers = []tr.SampleMethod{tr.CreateRgssSampler()}

	batchSize := 64
	var wg sync.WaitGroup

	for j := 0; j < opts.height; j += batchSize {
		for i := 0; i < opts.width; i += batchSize {
			wg.Add(1)
			go func(startX int, startY int, batch int) {
				defer wg.Done()
				t := sc.TraceParams{
					RayCast: g.Ray{Origin: g.Vector{X: 0, Y: 0, Z: 0}, Dir: g.Vector{X: 0, Y: 0, Z: 0}},
					Depth:   5,
					Value:   1.0,
				}
				for y := startY; y < startY+batch && y < opts.height; y++ {
					for x := startX; x < startX+batch && x < opts.width; x++ {
						targetSample := tr.SampleSingle(castOrigin, castCorner, dx, dy, x, y, samplers)
						color, _ := tr.TraceSample(targetSample, t, scene)
						testImage.Set(x, y, color)
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
		Size:     Size{Height: opts.width, Width: opts.height},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			CustomWidget{AssignTo: &mw.paintWidget, Paint: mw.renderImage},
		},
	}.Run()

	wg.Wait()

	testfile, _ := os.Create("test.png")
	png.Encode(testfile, testImage)
}
