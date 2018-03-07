package main

type light interface {
	// light calcluates the lighting values for a point in the scene
	// It returns the diffuse light and the specular light values
	// Note: ray.dir is not used
	// Accepts the incident, surface normal, and the scene
	light(vector, vector, scene) (ray, float64)
	getColor() floatColor
}
