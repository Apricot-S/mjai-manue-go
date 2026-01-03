package tile_test

import (
	"sort"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestSortTiles(t *testing.T) {
	names := [...]string{
		"?",
		"5sr", "5pr", "5mr",
		"C", "F", "P", "N", "W", "S", "E",
		"9s", "8s", "7s", "6s", "5s", "4s", "3s", "2s", "1s",
		"9p", "8p", "7p", "6p", "5p", "4p", "3p", "2p", "1p",
		"9m", "8m", "7m", "6m", "5m", "4m", "3m", "2m", "1m",
	}

	tiles := make(tile.Tiles, 0, len(names))
	for _, name := range names {
		t := tile.MustTileFromCode(name)
		tiles = append(tiles, *t)
	}
	sort.Sort(tiles)

	sortedNames := [...]string{
		"1m", "2m", "3m", "4m", "5m", "5mr", "6m", "7m", "8m", "9m",
		"1p", "2p", "3p", "4p", "5p", "5pr", "6p", "7p", "8p", "9p",
		"1s", "2s", "3s", "4s", "5s", "5sr", "6s", "7s", "8s", "9s",
		"E", "S", "W", "N", "P", "F", "C",
		"?",
	}

	for i, sortedName := range sortedNames {
		if tiles[i].String() != sortedName {
			t.Errorf("Expected %s but got %s", sortedName, tiles[i])
		}
	}
}
