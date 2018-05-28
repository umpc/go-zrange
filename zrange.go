package zrange

import "math"

const (
	earthSemiMajorAxis = 6378.137
	earthEquator       = math.Pi * earthSemiMajorAxis
)

var radiusToBits = precalcRadiusToBits()

func precalcRadiusToBits() []float64 {
	var radiusToBits []float64

	for bits, prevRadialBound := uint(4), earthEquator; bits < 64; bits += 2 {
		radiusToBits = append(radiusToBits, prevRadialBound/2)
		prevRadialBound = radiusToBits[len(radiusToBits)-1]
	}

	return radiusToBits
}
