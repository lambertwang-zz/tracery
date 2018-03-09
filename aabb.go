package main

import (
	"fmt"
	"math"
)

type slab struct {
	dNear, dFar float64
	normal      vector
}

// AKA an axis aligned bounding box
// The mean is used to partition the volumes by proximity
type aabb struct {
	slabs [3]slab
	mean  vector
	// min, max, mean vector
}

func (v aabb) clamp(sceneBounds aabb) aabb {
	for i := 0; i < 3; i++ {
		v.slabs[i].dNear = math.Max(v.slabs[i].dNear, sceneBounds.slabs[i].dNear)
		v.slabs[i].dFar = math.Min(v.slabs[i].dFar, sceneBounds.slabs[i].dFar)
	}
	return v
}

func combineAabb(sceneBounds aabb, boxes ...aabb) aabb {
	var newSlabs [3]slab

	for i := 0; i < 3; i++ {
		newSlabs[i].dNear = math.MaxFloat64
		newSlabs[i].dFar = -math.MaxFloat64
		newSlabs[i].normal = sceneBounds.slabs[i].normal
	}

	meanSum := vector{0, 0, 0}

	for _, box := range boxes {
		for i := 0; i < 3; i++ {
			newSlabs[i].dNear = math.Min(newSlabs[i].dNear, box.slabs[i].dNear)
			newSlabs[i].dFar = math.Max(newSlabs[i].dFar, box.slabs[i].dFar)
		}
		meanSum = meanSum.plus(box.mean)
	}

	return aabb{newSlabs, meanSum.scale(1 / float64(len(boxes)))}
}

func (s slab) intersects(r ray) (tNear, tFar float64) {
	denom := dotProduct(r.dir, s.normal)
	nDotO := dotProduct(r.origin, s.normal)
	if denom == 0 {
		return -math.MaxFloat64, math.MaxFloat64
	}
	tNear = -(nDotO + s.dNear) / denom
	tFar = -(nDotO + s.dFar) / denom
	if denom > 0 {
		tNear, tFar = tFar, tNear
	}
	return
}

func (v aabb) toString() string {
	return fmt.Sprintf("Min: x: %f y: %f z: %f\nMax x: %f y: %f z: %f\n",
		v.slabs[0].dNear,
		v.slabs[1].dNear,
		v.slabs[2].dNear,
		v.slabs[0].dFar,
		v.slabs[1].dFar,
		v.slabs[2].dFar,
	)
}

func (v aabb) intersects(r ray) (isIntersecting bool, tNear, tFar float64) {
	// Compute the incident of the ray with any face of the volume
	// If that point lies inside the other two dimensions of the volume, then it is inside the volume.
	tNear = -math.MaxFloat64
	tFar = math.MaxFloat64
	isIntersecting = true
	for _, slab := range v.slabs {
		newNear, newFar := slab.intersects(r)
		tNear = math.Max(newNear, tNear)
		tFar = math.Min(newFar, tFar)
		if tFar < tNear {
			isIntersecting = false
			return
		}
	}

	return
}
