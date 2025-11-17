package meld

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type ConcealedKan struct {
	consumed [4]tile.Tile
}

func NewConcealedKan(consumed [4]tile.Tile) (*ConcealedKan, error) {
	panic("")
}

func (k *ConcealedKan) Taken() *tile.Tile {
	panic("")
}

func (k *ConcealedKan) Consumed() []tile.Tile {
	return k.consumed[:]
}

func (k *ConcealedKan) Target() int {
	panic("")
}

func (k *ConcealedKan) ToTiles() []tile.Tile {
	return k.consumed[:]
}

func (k *ConcealedKan) ToBlock() block.Block {
	// Red five is sorted after normal, so RemoveRed() is not necessary.
	return block.MustQuad(k.consumed[0])
}

func (k *ConcealedKan) ToString() string {
	panic("")
}
