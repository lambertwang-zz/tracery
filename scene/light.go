package scene

import (
	g "github.com/lambertwang/tracery/geometry"
)

// Light TODO: Comment
type Light interface {
	// light calcluates the lighting values for a point in the scene
	// It returns the diffuse light and the specular light values
	// Note: ray.dir is not used
	// Accepts the incident, surface normal, and the scene
	light(g.Vector, g.Vector, Scene) (g.Ray, float64)
	GetColor() g.FloatColor
}
