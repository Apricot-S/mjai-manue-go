package hand

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
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
	t := tile.MustTileFromCode("?")
	tiles := make([]tile.Tile, h.tileCount)
	for i := range tiles {
		tiles[i] = t
	}
	return tiles
}

func (h *InvisibleHand) Draw(tile tile.Tile) (*InvisibleHand, error) {
	if h.tileCount >= maxNumTilesInHand {
		return nil, fmt.Errorf("cannot draw tile: hand already has %d tiles", h.tileCount)
	}
	return &InvisibleHand{tileCount: h.tileCount + 1}, nil
}

func (h *InvisibleHand) Discard(tile tile.Tile) (*InvisibleHand, error) {
	if h.tileCount <= 0 {
		return nil, fmt.Errorf("cannot discard tile: hand is empty")
	}
	return &InvisibleHand{tileCount: h.tileCount - 1}, nil
}

func (h *InvisibleHand) Call(m meld.Meld) (*InvisibleHand, error) {
	numConsumed := 0

	switch m.(type) {
	case *meld.Chii, *meld.Pon:
		numConsumed = 2
	case *meld.CalledKan:
		numConsumed = 3
	case *meld.ConcealedKan:
		numConsumed = 4
	case *meld.PromotedKan:
		numConsumed = 1
	}

	if h.tileCount < numConsumed {
		return nil, fmt.Errorf("cannot call %T: need %d tiles but hand has only %d", m, numConsumed, h.tileCount)
	}

	return &InvisibleHand{tileCount: h.tileCount - numConsumed}, nil
}
