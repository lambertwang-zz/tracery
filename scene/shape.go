package scene

import g "github.com/lambertwang/tracery/geometry"

// Shape TODO: Comment
type Shape interface {
	getMaterial() Material

	shouldTest(g.Ray) bool
	intersect(g.Ray) float64
	// Returns the incident and the normal
	traceTo(float64, g.Ray) (g.Vector, g.Vector)
	sampleC(incident, normal, dir g.Vector, sc Scene) g.FloatColor
	// Returns the incident point, all reflected directions, the color of the
	// object at that point, the reflectance, and the normal
	reflect(incident, normal, dir g.Vector, sc Scene) g.Vector
	refract(incident, normal, dir g.Vector, sc Scene) g.Vector

	// Takes in the scene max size and returns the Aabb of the shape
	bounds(Aabb) Aabb
}
