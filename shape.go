package main

import "image/color"

type shape interface {
	intersect(ray) (float64, float64)
	reflect(ray, float64, scene) (ray, color.RGBA, float64)
}

type material struct {
	color       color.RGBA
	reflectance float64
}
