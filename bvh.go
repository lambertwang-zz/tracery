package main

import (
	"fmt"
	"sort"
)

// Bounding volume heirarchy

type bvhNode struct {
	s        *shape
	bounds   aabb
	lhs, rhs *bvhNode
}

type byXMean []bvhNode
type byYMean []bvhNode
type byZMean []bvhNode

func (s byXMean) Less(i, j int) bool {
	return s[i].bounds.mean.x < s[j].bounds.mean.x
}
func (s byYMean) Less(i, j int) bool {
	return s[i].bounds.mean.y < s[j].bounds.mean.y
}
func (s byZMean) Less(i, j int) bool {
	return s[i].bounds.mean.z < s[j].bounds.mean.z
}
func (s byXMean) Len() int      { return len(s) }
func (s byXMean) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byYMean) Len() int      { return len(s) }
func (s byYMean) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byZMean) Len() int      { return len(s) }
func (s byZMean) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (b bvhNode) toString() string {
	return fmt.Sprintf("Children %d \n", b.count())
}

func (b bvhNode) count() int {
	if b.lhs == nil {
		if b.rhs == nil {
			return 1
		}
		return 1 + b.rhs.count()
	}
	if b.rhs == nil {
		return 1 + b.lhs.count()
	}
	return 1 + b.rhs.count() + b.lhs.count()
}

func constructHeirarchy(shapes *[]shape, sceneBounds aabb) *bvhNode {
	bvhQueue := make([]bvhNode, len(*shapes))

	for i, s := range *shapes {
		bvhQueue[i] = bvhNode{
			&(*shapes)[i],
			s.bounds(sceneBounds).clamp(sceneBounds),
			nil, nil,
		}
	}
	return constructBvhHelper(bvhQueue, sceneBounds, 0)
}

func constructBvhHelper(bvhQueue []bvhNode, sceneBounds aabb, depth int) *bvhNode {
	if len(bvhQueue) <= 0 {
		panic("BVH Queue is empty!")
	}
	if len(bvhQueue) == 1 {
		return &bvhNode{
			bvhQueue[0].s,
			bvhQueue[0].bounds,
			nil, nil,
		}
	}
	if len(bvhQueue) == 2 {
		return &bvhNode{
			nil,
			combineAabb(sceneBounds, bvhQueue[0].bounds, bvhQueue[1].bounds),
			constructBvhHelper(bvhQueue[:1], sceneBounds, depth+1),
			constructBvhHelper(bvhQueue[1:], sceneBounds, depth+1),
		}
	}
	switch depth % 3 {
	case 0:
		sort.Sort(byXMean(bvhQueue))
		break
	case 1:
		sort.Sort(byYMean(bvhQueue))
		break
	case 2:
		sort.Sort(byZMean(bvhQueue))
		break
	}
	midIndex := len(bvhQueue) / 2
	left := constructBvhHelper(bvhQueue[:midIndex], sceneBounds, depth+1)
	right := constructBvhHelper(bvhQueue[midIndex:], sceneBounds, depth+1)
	return &bvhNode{
		nil,
		combineAabb(sceneBounds, left.bounds, right.bounds),
		left, right,
	}
}
