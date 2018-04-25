package scene

import (
	"fmt"
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// Slab TODO: Comment
type Slab struct {
	DNear, DFar float64
	Normal      g.Vector
}

// Aabb AKA an axis aligned bounding box
// The mean is used to partition the volumes by proximity
type Aabb struct {
	Slabs [3]Slab
	Mean  g.Vector
}

func (v Aabb) clamp(sceneBounds Aabb) Aabb {
	for i := 0; i < 3; i++ {
		v.Slabs[i].DNear = math.Max(v.Slabs[i].DNear, sceneBounds.Slabs[i].DNear)
		v.Slabs[i].DFar = math.Min(v.Slabs[i].DFar, sceneBounds.Slabs[i].DFar)
	}
	return v
}

func combineAabb(sceneBounds Aabb, boxes ...Aabb) Aabb {
	var newSlabs [3]Slab

	for i := 0; i < 3; i++ {
		newSlabs[i].DNear = math.MaxFloat64
		newSlabs[i].DFar = -math.MaxFloat64
		newSlabs[i].Normal = sceneBounds.Slabs[i].Normal
	}

	meanSum := g.Vector{}

	for _, box := range boxes {
		for i := 0; i < 3; i++ {
			newSlabs[i].DNear = math.Min(newSlabs[i].DNear, box.Slabs[i].DNear)
			newSlabs[i].DFar = math.Max(newSlabs[i].DFar, box.Slabs[i].DFar)
		}
		meanSum = meanSum.Plus(box.Mean)
	}

	return Aabb{newSlabs, meanSum.Scale(1 / float64(len(boxes)))}
}

func (s Slab) intersects(r g.Ray) (tNear, tFar float64) {
	denom := r.Dir.Dot(s.Normal)
	nDotO := r.Origin.Dot(s.Normal)
	if denom == 0 {
		return -math.MaxFloat64, math.MaxFloat64
	}
	tNear = -(nDotO + s.DNear) / denom
	tFar = -(nDotO + s.DFar) / denom
	if denom > 0 {
		tNear, tFar = tFar, tNear
	}
	return
}

func (v Aabb) toString() string {
	return fmt.Sprintf("Min: x: %f y: %f z: %f\nMax x: %f y: %f z: %f\n",
		v.Slabs[0].DNear,
		v.Slabs[1].DNear,
		v.Slabs[2].DNear,
		v.Slabs[0].DFar,
		v.Slabs[1].DFar,
		v.Slabs[2].DFar,
	)
}

func (v Aabb) intersects(r g.Ray) (isIntersecting bool, tNear, tFar float64) {
	// Compute the incident of the ray with any face of the volume
	// If that point lies inside the other two dimensions of the volume, then it is inside the volume.
	tNear = -math.MaxFloat64
	tFar = math.MaxFloat64
	isIntersecting = true
	for _, Slab := range v.Slabs {
		newNear, newFar := Slab.intersects(r)
		tNear = math.Max(newNear, tNear)
		tFar = math.Min(newFar, tFar)
		if tFar < tNear {
			isIntersecting = false
			return
		}
	}

	return
}
