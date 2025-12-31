package hand

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type InvisibleHand struct {
	tileCount int
}

func NewInvisibleHand(tiles []tile.Tile) (*InvisibleHand, error) {
	sum := len(tiles)
	if sum > maxNumTilesInHand {
		return nil, fmt.Errorf("hand cannot contain 15 or more tiles: %d", sum)
	}

	return &InvisibleHand{tileCount: sum}, nil
}

func MustInvisibleHand(tiles []tile.Tile) *InvisibleHand {
	h, err := NewInvisibleHand(tiles)
	if err != nil {
		panic(err)
	}
	return h
}

func (h *InvisibleHand) ToTiles() []tile.Tile {
	t := *tile.MustTileFromCode("?")
	tiles := make([]tile.Tile, h.tileCount)
	for i := range tiles {
		tiles[i] = t
	}
	return tiles
}

func (h *InvisibleHand) Draw(tile *tile.Tile) (Hand, error) {
	if h.tileCount >= maxNumTilesInHand {
		return nil, fmt.Errorf("cannot draw tile: hand already has %d tiles", h.tileCount)
	}
	return &InvisibleHand{tileCount: h.tileCount + 1}, nil
}

func (h *InvisibleHand) Discard(tile *tile.Tile) (Hand, error) {
	if h.tileCount <= 0 {
		return nil, fmt.Errorf("cannot discard tile: hand is empty")
	}
	return &InvisibleHand{tileCount: h.tileCount - 1}, nil
}

func (h *InvisibleHand) Call(m meld.Meld) (Hand, error) {
	switch m.(type) {
	case *meld.Chii:
		return &InvisibleHand{tileCount: h.tileCount - 2}, nil
	}
	panic("unreachable")
}
