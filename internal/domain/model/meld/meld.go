package meld

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Meld interface {
	Taken() *tile.Tile
	Consumed() []tile.Tile
	Target() int
	ToTiles() []tile.Tile
	ToBlock() block.Block
	ToString() string
}

func isValidTarget(target int) bool {
	return 0 <= target && target <= 3
}

type Pon struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	target   int
	tiles    []tile.Tile
}

func NewPon(taken tile.Tile, consumed [2]tile.Tile, target int) (*Pon, error) {
	if !isValidTarget(target) {
		return nil, fmt.Errorf("invalid target: %d", target)
	}

	tiles := tile.Tiles{taken, consumed[0], consumed[1]}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return t.IsUnknown() }) {
		return nil, fmt.Errorf("unknown tile cannot use for Pon")
	}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return !taken.HasSameSymbol(&t) }) {
		return nil, fmt.Errorf("mismatch taken: %+v, consumed: %+v", taken, consumed)
	}

	numRed := 0
	for _, t := range tiles {
		if t.IsRed() {
			numRed++
		}
	}
	if numRed > 1 {
		return nil, fmt.Errorf("cannot use 2 or more red fives for Pon; taken: %+v, consumed: %+v", taken, consumed)
	}

	sort.Sort(tiles)

	csm := tile.Tiles(consumed[:])
	sort.Sort(csm)

	return &Pon{
		taken:    taken,
		consumed: [2]tile.Tile(csm),
		target:   target,
		tiles:    tiles,
	}, nil
}

func (p *Pon) Taken() *tile.Tile {
	return &p.taken
}

func (p *Pon) Consumed() []tile.Tile {
	return p.consumed[:]
}

func (p *Pon) Target() int {
	return p.target
}

func (p *Pon) ToTiles() []tile.Tile {
	return p.tiles
}

func (p *Pon) ToBlock() block.Block {
	// Red five is sorted after normal, so RemoveRed() is not necessary.
	return block.MustTriplet(p.tiles[0])
}

func ToString() string {
	panic("")
}
