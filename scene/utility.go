package scene

import (
	"math"

	g "github.com/lambertwang/tracery/geometry"
)

// NearestIntersection TODO: Comment
func NearestIntersection(r g.Ray, s Scene) (nearestT float64, nearestShape *Shape) {
	nearestT, nearestShape = checkBvh(r, s.bvh)
	// Following code is before BVH
	/*
		nearestT = math.MaxFloat64
		nearestShape = nil
		for i, shape := range s.shapes {
			t := shape.intersect(r)
			if t > 0.001 && t < nearestT {
				nearestT = t
				nearestShape = &s.shapes[i]
			}
		}
	*/
	return
}

// Returns the nearestT and the shape intersected
func checkBvh(r g.Ray, n *bvhNode) (float64, *Shape) {
	intersectsBounds, _, _ := n.bounds.intersects(r)
	if intersectsBounds {
		if n.s == nil {
			lhsT, lhsS := checkBvh(r, n.lhs)
			rhsT, rhsS := checkBvh(r, n.rhs)
			if lhsS != nil && lhsT < rhsT {
				return lhsT, lhsS
			}
			if rhsS != nil {
				return rhsT, rhsS
			}
		} else {
			t := (*n.s).intersect(r)
			if t > 0.001 {
				return t, n.s
			}
		}
	}
	return math.MaxFloat64, nil
}
