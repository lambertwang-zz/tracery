package main

import (
	"image/color"
	"math"
	"strconv"
)

const reflectanceThreshold = 0.01

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func vecStr(v vector) string {
	return "x = " + floatStr(v.x) +
		", y = " + floatStr(v.y) +
		", z = " + floatStr(v.z)
}

func meanColor(colors ...color.RGBA) color.RGBA {
	var r, g, b, a int

	for _, c := range colors {
		r += int(c.R)
		g += int(c.G)
		b += int(c.B)
		a += int(c.A)
	}
	r /= len(colors)
	g /= len(colors)
	b /= len(colors)
	a /= len(colors)

	return color.RGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		// uint8(a),
		255,
	}
}

func scaleColor(c color.RGBA, f float64) color.RGBA {
	scale := math.Max(0.0, f)
	return color.RGBA{
		uint8(math.Min(float64(c.R)*scale, 255)),
		uint8(math.Min(float64(c.G)*scale, 255)),
		uint8(math.Min(float64(c.B)*scale, 255)),
		c.A,
	}
}

func clamp16to8(x uint16) uint8 {
	if x > 255 {
		return 255
	}
	return uint8(x)
}

func addColor(lhs color.RGBA, rhs color.RGBA) color.RGBA {
	r := clamp16to8(uint16(lhs.R) + uint16(rhs.R))
	g := clamp16to8(uint16(lhs.G) + uint16(rhs.G))
	b := clamp16to8(uint16(lhs.B) + uint16(rhs.B))
	a := clamp16to8(uint16(lhs.A) + uint16(rhs.A))
	return color.RGBA{r, g, b, a}
}

const (
	uniformSampling = 1
	rgssSampling    = 2
)

func castQuad(origin vector, corner vector, dx vector, dy vector, stepsX int, stepsY int) ([]ray, int) {
	rayCasts := make([]ray, stepsX*stepsY)

	for y := 0; y < stepsY; y++ {
		for x := 0; x < stepsX; x++ {
			rayCasts[y*stepsX+x] = lineToRay(
				origin,
				addVectors(
					corner,
					dx.scale(float64(x)),
					dy.scale(float64(y)),
				),
			)
		}
	}

	return rayCasts, 1
}

func castQuadRgss(origin vector, corner vector, dx vector, dy vector, stepsX int, stepsY int) ([]ray, int) {
	rayCasts := make([]ray, stepsX*stepsY*4)

	for y := 0; y < stepsY; y++ {
		for x := 0; x < stepsX; x++ {
			rayCasts[(y*stepsX+x)*4] = lineToRay(
				origin,
				addVectors(
					corner,
					dx.scale(float64(x)+.125),
					dy.scale(float64(y)+.375),
				),
			)
			rayCasts[(y*stepsX+x)*4+1] = lineToRay(
				origin,
				addVectors(
					corner,
					dx.scale(float64(x)+.625),
					dy.scale(float64(y)+.125),
				),
			)
			rayCasts[(y*stepsX+x)*4+2] = lineToRay(
				origin,
				addVectors(
					corner,
					dx.scale(float64(x)+.375),
					dy.scale(float64(y)+.875),
				),
			)
			rayCasts[(y*stepsX+x)*4+3] = lineToRay(
				origin,
				addVectors(
					corner,
					dx.scale(float64(x)+.875),
					dy.scale(float64(y)+.625),
				),
			)
		}
	}

	return rayCasts, 4
}

type traceParams struct {
	reflections int
	reflectance float64
	color       color.RGBA
}

func trace(ray ray, t traceParams, s scene) color.RGBA {
	// Compute intersections
	nearestT := math.MaxFloat64
	index := -1
	for i, shape := range s.shapes {
		t0, t1 := shape.intersect(ray)
		// Assume ray cannot cast from inside of the sphere
		if t0 > 0.01 && t1 > 0.01 {
			closeT := math.Min(t0, t1)
			if closeT < nearestT {
				nearestT = closeT
				index = i
			}
		}
	}
	if index < 0 {
		return t.color
	}

	if t.reflections <= 0 {
		return t.color
	}

	reflection, c, reflectance := s.shapes[index].reflect(ray, nearestT, s)

	t.reflections--
	t.color = addColor(t.color, scaleColor(c, t.reflectance))
	t.reflectance *= reflectance
	if reflectance < reflectanceThreshold {
		return t.color
	}

	return trace(reflection, t, s)
}
