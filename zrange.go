// Package zrange implements an efficient algorithm for performing geospatial
// range queries with Geohash-encoded keys and a search radius.
//
// The RadialRange method appears to be sufficient for range queries of around
// 5,000km or less. Changes that efficiently add support for larger query ranges
// are welcome.
//
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
