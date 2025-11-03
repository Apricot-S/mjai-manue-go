package hand

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

const maxNumTilesInHand = 14
const maxCopies = 4

type Hand struct {
	tileCounts [tile.NumTileType38]int
}

func NewHand(tiles []tile.Tile) (*Hand, error) {
	h := Hand{}
	for _, t := range tiles {
		id := t.ID()
		h.tileCounts[id]++

		if h.tileCounts[id] > maxCopies {
			return nil, fmt.Errorf("tiles contains five identical tiles: %s", t.Code())
		}
	}
	return &h, nil
}

func (h *Hand) ToTiles() []tile.Tile {
	tiles := make([]tile.Tile, 0, maxNumTilesInHand)
	for i, c := range h.tileCounts {
		for range c {
			tiles = append(tiles, *tile.MustTileFromID(i))
		}
	}
	return tiles
}
