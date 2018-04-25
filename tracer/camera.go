package tracer

import (
	g "github.com/lambertwang/tracery/geometry"
)

// Camera TODO: Comment
type Camera struct {
	origin         g.Vector
	dx, dy         g.Vector
	stepsX, stepsY int
}

func createCamera(origin, target g.Vector) (outCam Camera) {
	return Camera{
		origin,
		target, target,
		0, 0,
	}
}
