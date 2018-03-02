package main

type light interface {
	light(ray, vector, scene) float64
}
