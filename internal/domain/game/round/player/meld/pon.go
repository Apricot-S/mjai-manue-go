package meld

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Pon struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	target   seat.Seat
	tiles    []tile.Tile
}

func NewPon(taken tile.Tile, consumed [2]tile.Tile, target seat.Seat) (*Pon, error) {
	tiles := tile.Tiles{taken, consumed[0], consumed[1]}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return t.IsUnknown() }) {
		return nil, fmt.Errorf("unknown tile cannot use for Pon")
	}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return !taken.HasSameSymbol(&t) }) {
		return nil, fmt.Errorf("mismatch taken: %+v, consumed: %+v", taken, consumed)
	}
	if taken.IsSuits() && taken.Number() == 5 && countRed(tiles) > 1 {
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

func MustPon(taken tile.Tile, consumed [2]tile.Tile, target seat.Seat) *Pon {
	p, err := NewPon(taken, consumed, target)
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Pon) Taken() *tile.Tile {
	return &p.taken
}

func (p *Pon) Consumed() []tile.Tile {
	return p.consumed[:]
}

func (p *Pon) Target() *seat.Seat {
	return &p.target
}

func (p *Pon) ToTiles() []tile.Tile {
	return p.tiles
}

func (p Pon) String() string {
	return meldToString(&p)
}

func (p *Pon) SwapCallTiles() []tile.Tile {
	// Red five is sorted after normal, so RemoveRed() is not necessary.
	if p.taken.IsSuits() && p.taken.Number() == 5 {
		return []tile.Tile{p.tiles[0], p.tiles[0].AddRed()}
	}
	return []tile.Tile{p.tiles[0]}
}
