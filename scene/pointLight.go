package scene

import (
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// PointLight TODO: Comment
type PointLight struct {
	Center g.Vector
	Radius float64
	Color  g.FloatColor
}

func (l PointLight) light(incident g.Vector, normal g.Vector, s Scene) (g.Ray, float64) {
	// Compute intersections
	toRay := g.LineToRay(incident, l.Center)
	nearestT, _ := NearestIntersection(toRay, s)

	v := l.Center.Minus(incident)

	if v.Magnitude() > l.Center.Minus(toRay.Incident(nearestT)).Magnitude() {
		return toRay, 0.0
	}

	distance := incident.Minus(l.Center).Magnitude()
	attenuation := 1.0 / (1.0 + (2.0/l.Radius)*distance + (1.0/math.Pow(l.Radius, 2))*math.Pow(distance, 2))

	return toRay, attenuation
}

// GetColor TODO: Comment
func (l PointLight) GetColor() g.FloatColor {
	return l.Color
}
