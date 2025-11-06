package meld

import (
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
