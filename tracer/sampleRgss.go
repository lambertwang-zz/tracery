package tracer

import (
	g "github.com/lambertwang/tracery/geometry"
)

// CreateRgssSampler TODO: Comment
func CreateRgssSampler() SampleMethod {
	return func(
		origin g.Vector, target g.Vector,
		dx g.Vector, dy g.Vector,
		rays *[]g.Ray) {
		*rays = append(
			*rays,
			g.Ray{
				Origin: origin,
				Dir: g.AddVectors(target,
					dx.Scale(.125),
					dy.Scale(.375),
				),
			},
			g.Ray{
				Origin: origin,
				Dir: g.AddVectors(target,
					dx.Scale(.625),
					dy.Scale(.125),
				),
			},
			g.Ray{
				Origin: origin,
				Dir: g.AddVectors(target,
					dx.Scale(.375),
					dy.Scale(.875),
				),
			},
			g.Ray{
				Origin: origin,
				Dir: g.AddVectors(target,
					dx.Scale(.875),
					dy.Scale(.625),
				),
			},
		)
	}
}
