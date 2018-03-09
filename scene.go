package main

type scene struct {
	shapes       []shape
	lights       []light
	ambientLight float64
	bvh          *bvhNode
}
