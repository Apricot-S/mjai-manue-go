package hand

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"

// TileCounts34 is the 34-kind tile count representation. Red fives are merged
// into the corresponding normal five counters before values reach this type.
type TileCounts34 [tile.NumTileType34]int

func (tc34 *TileCounts34) ToTiles() []tile.Tile {
	tiles := make([]tile.Tile, 0, tc34.NumTiles())
	for i, c := range tc34 {
		for range c {
			tiles = append(tiles, tile.MustTileFromID(i))
		}
	}
	return tiles
}

func (tc34 *TileCounts34) NumTiles() int {
	sum := 0
	for _, c := range tc34 {
		sum += c
	}
	return sum
}
