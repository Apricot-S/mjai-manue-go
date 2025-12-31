package hand

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tilecount"
)

type VisibleHand struct {
	tileCounts [tile.NumTileType37]int
	numTiles   int
}

func NewVisibleHand(tiles []tile.Tile) (*VisibleHand, error) {
	tileCounts := [tile.NumTileType37]int{}
	sum := 0

	for _, t := range tiles {
		id := t.ID()
		if id >= tile.NumTileType37 {
			return nil, fmt.Errorf("visible hand cannot contain unknown tiles")
		}

		tileCounts[id]++
		sum++
		if tileCounts[id] > maxCopies {
			return nil, fmt.Errorf("hand cannot contain five identical tiles: %s", t)
		}
		if t.IsRed() && tileCounts[id] > 1 {
			return nil, fmt.Errorf("hand cannot contain multiple red fives of the same suit: %s", t)
		}
	}

	if sum > maxNumTilesInHand {
		return nil, fmt.Errorf("hand cannot contain 15 or more tiles: %d", sum)
	}

	return &VisibleHand{tileCounts: tileCounts, numTiles: sum}, nil
}

func MustVisibleHand(tiles []tile.Tile) *VisibleHand {
	h, err := NewVisibleHand(tiles)
	if err != nil {
		panic(err)
	}
	return h
}

func (h *VisibleHand) ToTiles() []tile.Tile {
	tiles := make([]tile.Tile, 0, maxNumTilesInHand)
	for i, c := range h.tileCounts {
		for range c {
			tiles = append(tiles, *tile.MustTileFromID(i))
		}
	}
	return tiles
}

func (h *VisibleHand) ToTileCounts34() *tilecount.TileCounts34 {
	tc := tilecount.TileCounts34(h.tileCounts[:34])
	tc[4] += h.tileCounts[34]
	tc[13] += h.tileCounts[35]
	tc[22] += h.tileCounts[36]
	return &tc
}

func (h *VisibleHand) Draw(tile *tile.Tile) (Hand, error) {
	if tile.IsUnknown() {
		return nil, fmt.Errorf("visible hand cannot draw an unknown tile")
	}

	if h.numTiles >= maxNumTilesInHand {
		return nil, fmt.Errorf("cannot draw tile: hand already has %d tiles", h.numTiles)
	}

	id := tile.ID()
	tileCounts := h.tileCounts
	if tileCounts[id] >= maxCopies {
		return nil, fmt.Errorf("cannot draw tile: hand already has four identical tiles: %s", tile)
	}
	if tile.IsRed() && tileCounts[id] >= 1 {
		return nil, fmt.Errorf("cannot draw tile: hand already has a red five: %s", tile)
	}

	tileCounts[id]++
	return &VisibleHand{tileCounts: tileCounts, numTiles: h.numTiles + 1}, nil
}

func (h *VisibleHand) Discard(tile *tile.Tile) (Hand, error) {
	if tile.IsUnknown() {
		return nil, fmt.Errorf("visible hand cannot discard an unknown tile")
	}

	id := tile.ID()
	tileCounts := h.tileCounts
	if tileCounts[id] <= 0 {
		return nil, fmt.Errorf("cannot discard tile: %s is not in the hand", tile)
	}

	tileCounts[id]--
	return &VisibleHand{tileCounts: tileCounts, numTiles: h.numTiles - 1}, nil
}

func (h *VisibleHand) Call(m meld.Meld) (Hand, error) {
	panic("")
}
