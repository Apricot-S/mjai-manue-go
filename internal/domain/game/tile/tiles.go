package tile

import (
	"slices"
)

var tileSortKeyTable = [NumTileType38]int{
	0, 1, 2, 3, 4, 6, 7, 8, 9, // m
	10, 11, 12, 13, 14, 16, 17, 18, 19, // p
	20, 21, 22, 23, 24, 26, 27, 28, 29, // s
	30, 31, 32, 33, 34, 35, 36, // z
	5, 15, 25, // red
	37, // unknown
}

type Tiles []Tile

func (t *Tile) sortKey() int {
	return tileSortKeyTable[t.ID()]
}

func (t *Tile) compareTo(other *Tile) int {
	if other == nil {
		panic("Other tile is nil")
	}
	return t.sortKey() - other.sortKey()
}

func (ts Tiles) Sort() {
	slices.SortFunc(ts, func(a, b Tile) int {
		return a.compareTo(&b)
	})
}

func (ts *Tiles) Distinct(exclude func(Tile) bool) Tiles {
	ret := slices.Clone(*ts)
	ret.Sort()

	ret = slices.CompactFunc(ret, func(a, b Tile) bool {
		return a.ID() == b.ID()
	})

	if exclude == nil {
		return ret
	}

	return slices.DeleteFunc(ret, exclude)
}
