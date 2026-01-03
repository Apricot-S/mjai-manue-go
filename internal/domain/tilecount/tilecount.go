package tilecount

import "github.com/Apricot-S/mjai-manue-go/internal/domain/tile"

type TileCounts34 [tile.NumTileType34]int

func (tc34 *TileCounts34) ToTiles() []tile.Tile {
	tiles := make([]tile.Tile, 0)
	for i, c := range tc34 {
		for range c {
			tiles = append(tiles, *tile.MustTileFromID(i))
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
