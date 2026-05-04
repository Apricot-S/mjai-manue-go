package tile_test

import (
	"reflect"
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
	tiles.Sort()

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

func TestTiles_ContainsUnknown(t *testing.T) {
	tests := []struct {
		name  string
		tiles tile.Tiles
		want  bool
	}{
		{
			name:  "nil tiles",
			tiles: nil,
			want:  false,
		},
		{
			name:  "empty tiles",
			tiles: tile.Tiles{},
			want:  false,
		},
		{
			name: "no unknown",
			tiles: tile.Tiles{
				*tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("5mr"),
				*tile.MustTileFromCode("E"),
			},
			want: false,
		},
		{
			name: "contains unknown",
			tiles: tile.Tiles{
				*tile.MustTileFromCode("1m"),
				*tile.MustTileFromCode("?"),
				*tile.MustTileFromCode("E"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tiles.ContainsUnknown()
			if got != tt.want {
				t.Errorf("ContainsUnknown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTiles_Distinct(t *testing.T) {
	tests := []struct {
		name    string
		tiles   tile.Tiles
		exclude func(tile.Tile) bool
		want    tile.Tiles
	}{
		{
			name:    "nil tiles",
			tiles:   nil,
			exclude: nil,
			want:    nil,
		},
		{
			name:    "sort tiles",
			tiles:   tile.Tiles{*tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m")},
			exclude: nil,
			want:    tile.Tiles{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
		{
			name:    "remove duplicate tiles",
			tiles:   tile.Tiles{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m")},
			exclude: nil,
			want:    tile.Tiles{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr")},
		},
		{
			name:    "exclude tiles",
			tiles:   tile.Tiles{*tile.MustTileFromCode("5m"), *tile.MustTileFromCode("5mr"), *tile.MustTileFromCode("5m")},
			exclude: func(t tile.Tile) bool { return t.IsRed() },
			want:    tile.Tiles{*tile.MustTileFromCode("5m")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tiles.Distinct(tt.exclude)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Distinct() = %v, want %v", got, tt.want)
			}
		})
	}
}
