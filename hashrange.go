package zrange

import "sort"

// HashRange contains a minimum and maximum Geohash integer range value.
type HashRange struct {
	Min, Max uint64
}

// HashRanges is a list of ranges containing minimum and maximum Geohash integers,
// used for performing range queries.
type HashRanges []HashRange

// SortMinAsc sorts a list of hash ranges in ascending order by their items'
// Min fields.
func (hashRangeList HashRanges) SortMinAsc() HashRanges {
	sort.Sort(hashRangesMinAscSorter(hashRangeList))
	return hashRangeList
}

// CombineRanges merges each overlapping range.
// The input list of hash ranges are expected to be sorted in ascending order by
// their items' Min fields.
func (hashRangeList HashRanges) CombineRanges() HashRanges {
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
