package meld

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Pon struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	target   playerid.PlayerID
	tiles    []tile.Tile
}

func NewPon(taken tile.Tile, consumed [2]tile.Tile, target playerid.PlayerID) (*Pon, error) {
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

func (p *Pon) Taken() *tile.Tile {
	return &p.taken
}

func (p *Pon) Consumed() []tile.Tile {
	return p.consumed[:]
}

func (p *Pon) Target() *playerid.PlayerID {
	return &p.target
}

func (p *Pon) ToTiles() []tile.Tile {
	return p.tiles
}

func (p *Pon) ToBlock() block.Block {
	// Red five is sorted after normal, so RemoveRed() is not necessary.
	return block.MustTriplet(p.tiles[0])
}

func (p *Pon) String() string {
	return meldToString(p)
}

func (p *Pon) SwapCallTiles() []tile.Tile {
	// Red five is sorted after normal, so RemoveRed() is not necessary.
	if p.taken.IsSuits() && p.taken.Number() == 5 {
		return []tile.Tile{p.tiles[0], *p.tiles[0].AddRed()}
	}
	return []tile.Tile{p.tiles[0]}
}
