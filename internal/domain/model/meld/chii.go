package meld

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type chiiType int

const (
	chiiTypeLow chiiType = iota + 1
	chiiTypeMiddle
	chiiTypeHigh
)

type Chii struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	target   playerid.PlayerID
	tiles    []tile.Tile
	ty       chiiType
}

func NewChii(taken tile.Tile, consumed [2]tile.Tile, target playerid.PlayerID) (*Chii, error) {
	tiles := tile.Tiles{taken, consumed[0], consumed[1]}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return !t.IsSuits() }) {
		return nil, fmt.Errorf("honors or unknown tile cannot use for Chii; taken: %+v, consumed: %+v", taken, consumed)
	}

	sort.Sort(tiles)
	if tiles[0].Number() > 7 {
		return nil, fmt.Errorf("Chii cannot start with 8 or 9; taken: %+v, consumed: %+v", taken, consumed)
	}
	if !tiles[0].Next(1).HasSameSymbol(&tiles[1]) || !tiles[0].Next(2).HasSameSymbol(&tiles[2]) {
		return nil, fmt.Errorf("")
	}

	csm := tile.Tiles(consumed[:])
	sort.Sort(csm)

	var ty chiiType
	switch {
	case taken.HasSameSymbol(&tiles[0]):
		ty = chiiTypeLow
	case taken.HasSameSymbol(&tiles[1]):
		ty = chiiTypeMiddle
	case taken.HasSameSymbol(&tiles[2]):
		ty = chiiTypeHigh
	}

	return &Chii{
		taken:    taken,
		consumed: [2]tile.Tile(csm),
		target:   target,
		tiles:    tiles,
		ty:       ty,
	}, nil
}

func (c *Chii) Taken() *tile.Tile {
	return &c.taken
}

func (c *Chii) Consumed() []tile.Tile {
	return c.consumed[:]
}

func (c *Chii) Target() *playerid.PlayerID {
	return &c.target
}

func (c *Chii) ToTiles() []tile.Tile {
	return c.tiles
}

func (c *Chii) ToBlock() block.Block {
	return block.MustSequence(*c.tiles[0].RemoveRed())
}

func (c *Chii) String() string {
	return meldToString(c)
}

func (c *Chii) SwapCallTiles() []tile.Tile {
	n := c.taken.Number()

	switch c.ty {
	case chiiTypeLow:
		switch n {
		case 7:
			return []tile.Tile{c.taken}
		case 5:
			return []tile.Tile{*c.taken.RemoveRed(), *c.taken.AddRed(), *c.taken.Next(3)}
		case 2:
			next := c.taken.Next(3)
			return []tile.Tile{c.taken, *next, *next.AddRed()}
		default:
			return []tile.Tile{c.taken, *c.taken.Next(3)}
		}
	case chiiTypeMiddle:
		switch n {
		case 5:
			return []tile.Tile{*c.taken.RemoveRed(), *c.taken.AddRed()}
		default:
			return []tile.Tile{c.taken}
		}
	case chiiTypeHigh:
		switch n {
		case 3:
			return []tile.Tile{c.taken}
		case 5:
			return []tile.Tile{*c.taken.Next(-3), *c.taken.RemoveRed(), *c.taken.AddRed()}
		case 8:
			prev := c.taken.Next(-3)
			return []tile.Tile{*prev, *prev.AddRed(), c.taken}
		default:
			return []tile.Tile{*c.taken.Next(-3), c.taken}
		}
	}

	panic("unreachable: Chii.SwapCallTiles()")
}
