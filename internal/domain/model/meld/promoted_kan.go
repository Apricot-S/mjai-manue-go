package meld

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type PromotedKan struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	added    tile.Tile
	target   int
	tiles    []tile.Tile
}

func NewPromotedKan(taken tile.Tile, consumed [2]tile.Tile, added tile.Tile, target int) (*PromotedKan, error) {
	if !isValidTarget(target) {
		return nil, fmt.Errorf("invalid target: %d", target)
	}

	return &PromotedKan{
		taken:    taken,
		consumed: consumed,
		added:    added,
		target:   target,
		tiles:    []tile.Tile{},
	}, nil
}

func (k *PromotedKan) Taken() *tile.Tile {
	return &k.taken
}

func (k *PromotedKan) Consumed() []tile.Tile {
	return k.consumed[:]
}

func (k *PromotedKan) Added() *tile.Tile {
	return &k.added
}

func (k *PromotedKan) Target() int {
	return k.target
}
