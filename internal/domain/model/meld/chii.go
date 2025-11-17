package meld

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Chii struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	target   int
	tiles    []tile.Tile
}

func NewChii(taken tile.Tile, consumed [2]tile.Tile, target int) (*Chii, error) {
	if !isValidTarget(target) {
		return nil, fmt.Errorf("invalid target: %d", target)
	}

	tiles := tile.Tiles{taken, consumed[0], consumed[1]}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return !t.IsSuits() }) {
		return nil, fmt.Errorf("honors or unknown tile cannot use for Chii")
	}

	csm := tile.Tiles(consumed[:])

	return &Chii{
		taken:    taken,
		consumed: [2]tile.Tile(csm),
		target:   target,
		tiles:    tiles,
	}, nil
}

func (c *Chii) Taken() *tile.Tile {
	return &c.taken
}

func (c *Chii) Consumed() []tile.Tile {
	return c.consumed[:]
}

func (c *Chii) Target() int {
	return c.target
}

func (c *Chii) ToTiles() []tile.Tile {
	return c.tiles
}

func (c *Chii) ToBlock() block.Block {
	return block.MustSequence(*c.tiles[0].RemoveRed())
}

func (c *Chii) ToString() string {
	return meldToString(c)
}
