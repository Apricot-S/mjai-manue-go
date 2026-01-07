package tile

import (
	"log"
	"slices"
	"sort"
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
		log.Panic("Other tile is nil")
	}
	return t.sortKey() - other.sortKey()
}

func (ts Tiles) Len() int {
	return len(ts)
}

func (ts Tiles) Less(i, j int) bool {
	return ts[i].compareTo(&ts[j]) < 0
}

func (ts Tiles) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

func (ts *Tiles) Distinct(exclude func(Tile) bool) Tiles {
	ret := slices.Clone(*ts)
	sort.Sort(ret)

	ret = slices.CompactFunc(ret, func(a, b Tile) bool {
		return a.ID() == b.ID()
	})

	if exclude == nil {
		return ret
	}

	ret = slices.DeleteFunc(ret, func(t Tile) bool {
		return exclude(t)
	})

	return ret
}
