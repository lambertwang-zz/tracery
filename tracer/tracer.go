package tracer

import (
	"image/color"

	g "github.com/lambertwang/tracery/geometry"
	sc "github.com/lambertwang/tracery/scene"
)

const reflectanceThreshold = 0.01

// TraceSample TODO: comment
func TraceSample(sa Sample, t sc.TraceParams, s sc.Scene) (color.RGBA, float64) {
	c := make([]g.FloatColor, len(sa.dispatchedRays))

	outDepth := -1.0
	for i, r := range sa.dispatchedRays {
		newParams := sc.TraceParams{
			RayCast: r, Depth: t.Depth, Value: t.Value,
		}
		c[i], outDepth = s.Trace(newParams)
	}
	return g.MeanColor(c).ToRgba(), outDepth
}
