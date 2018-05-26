package zrange

import (
	"math"
	"sort"

	"github.com/mmcloughlin/geohash"
)

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

// RadialRange uses a radius in kilometers, a latitude, and a longitude to return
// a slice of one or more ranges of keys that can be used to efficiently perform
// geohash-based spatial queries.
//
// This method uses an algorithm that was derived from the "Search" section of this page:
// https://web.archive.org/web/20180526044934/https://github.com/yinqiwen/ardb/wiki/Spatial-Index#search
//
// RadialRange expands upon the ideas referenced above, by:
//
// • Sorting key ranges
//
// • Combining overlapping key ranges
//
// • Handling overflows resulting from bitshifting, such as when querying for: (-90, -180)
//
func RadialRange(params RadialRangeParams) HashRanges {
	return params.
		setDefaults().
		findNeighboringRanges().
		combineRanges()
}

// RadialRangeParams defaults to expecting 64-bit geohash-encoded keys.
type RadialRangeParams struct {
	BitsOfPrecision uint
	Radius,
	Latitude,
	Longitude float64
}

func (params RadialRangeParams) setDefaults() RadialRangeParams {
	if params.BitsOfPrecision == 0 {
		params.BitsOfPrecision = 64
	}
	return params
}

func (params RadialRangeParams) radiusToBits() uint {
	const initialSignificantBits = 2

	for i := len(radiusToBits) - 1; i > 0; i-- {
		if params.Radius < radiusToBits[i] {
			return uint(i*2 + initialSignificantBits)
		}
	}

	return uint(initialSignificantBits)
}

func (params RadialRangeParams) findNeighboringRanges() HashRanges {
	rangeBits := params.radiusToBits()

	queryPoint := geohash.EncodeIntWithPrecision(
		params.Latitude,
		params.Longitude,
		rangeBits,
	)

	neighborList := neighbors(geohash.NeighborsIntWithPrecision(queryPoint, rangeBits))
	neighborList = append(neighborList, queryPoint)

	rangeBitsDiff := params.BitsOfPrecision - rangeBits
	return neighborList.expandRanges(rangeBitsDiff)
}

type neighbors []uint64

func (neighborList neighbors) expandRanges(rangeBitsDiff uint) HashRanges {
	hashRangeList := make(HashRanges, 0, len(neighborList))

	for _, neighbor := range neighborList {
		min := neighbor << rangeBitsDiff
		max := (neighbor + 1) << rangeBitsDiff

		// Handle overflows near the outer edges.
		// Example: (-90.0, -180.0)
		if min > max {
			continue
		}

		hashRangeList = append(hashRangeList, HashRange{
			Min: min,
			Max: max,
		})
	}

	return hashRangeList
}

type hashRangesMinAscSorter HashRanges

func (s hashRangesMinAscSorter) Len() int {
	return len(s)
}
func (s hashRangesMinAscSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s hashRangesMinAscSorter) Less(i, j int) bool {
	return s[i].Min < s[j].Min
}

// HashRange contains a minimum and maximum geohash integer range value.
type HashRange struct {
	Min, Max uint64
}

// HashRanges is a list of ranges containing minimum and maximum geohash integers,
// used for performing range queries.
type HashRanges []HashRange

func (hashRangeList HashRanges) combineRanges() HashRanges {
	sort.Sort(hashRangesMinAscSorter(hashRangeList))

	combinedHashRangeList := hashRangeList[:0]
	for i := 0; i < len(hashRangeList)-1; i++ {
		hashRange := hashRangeList[i]
		nextHashRange := hashRangeList[i+1]

		if hashRange.Max == nextHashRange.Min {
			hashRange.Max = nextHashRange.Max
		}

		if hashRange.Max == nextHashRange.Max {
			hashRangeList[i+1].Min = hashRange.Min
			continue
		}

		combinedHashRangeList = append(combinedHashRangeList, hashRange)
	}

	return append(combinedHashRangeList, hashRangeList[len(hashRangeList)-1])
}
