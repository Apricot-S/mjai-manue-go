package meld

import (
	"fmt"

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

	panic("")
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
