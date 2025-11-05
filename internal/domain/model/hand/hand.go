package hand

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

const maxNumTilesInHand = 14
const maxCopies = 4

type Hand struct {
	tileCounts [tile.NumTileType38]int
}

func NewHand(tiles []tile.Tile) (*Hand, error) {
	tileCounts := [tile.NumTileType38]int{}
	for _, t := range tiles {
		id := t.ID()
		tileCounts[id]++

		if id >= tile.NumTileType37 {
			// There can be any number of unknowns.
			continue
		}
		if tileCounts[id] > maxCopies {
			return nil, fmt.Errorf("tiles cannot contain five identical tiles: %s", t.Code())
		}
		if t.IsRed() && tileCounts[id] > 1 {
			return nil, fmt.Errorf("tiles cannot contain multiple red fives of the same suit: %s", t.Code())
		}
		if t.IsSuits() && t.Number() == 5 && tileCounts[id] > 3 {
			return nil, fmt.Errorf("tiles cannot contain four normal fives of the same suit: %s", t.Code())
		}
	}

	sum := 0
	for _, c := range tileCounts {
		sum += c
	}
	if sum > maxNumTilesInHand {
		return nil, fmt.Errorf("tiles cannot contain 15 or more tiles: %d", sum)
	}

	return &Hand{tileCounts: tileCounts}, nil
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

func (h *Hand) ToTileCounts34() *tilecount.TileCounts34 {
	tc := tilecount.TileCounts34(h.tileCounts[:34])
	tc[4] += h.tileCounts[34]
	tc[13] += h.tileCounts[35]
	tc[22] += h.tileCounts[36]
	return &tc
}
