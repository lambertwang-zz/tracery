package main

type light interface {
	// light calcluates the lighting values for a point in the scene
	// It returns the diffuse light and the specular light values
	light(ray, vector, scene) (ray, float64)
}
