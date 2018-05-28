package zrange

type neighbors []uint64

func (neighborList neighbors) shiftIntoRanges(rangeBitsDiff uint) HashRanges {
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
