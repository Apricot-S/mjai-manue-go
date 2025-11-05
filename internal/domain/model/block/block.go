package block

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type BlockType int

// const (
// 	Sequence BlockType = iota + 1
// 	Triplet
// 	Quad
// 	Pair
// )

type Block interface {
	ToTiles() []tile.Tile
}

type Pair struct {
	tiles []tile.Tile
}

func NewPair(t tile.Tile) (*Pair, error) {
	if t.IsUnknown() {
		return nil, fmt.Errorf("cannot create pair from unknown tile")
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create pair from red five")
	}
	return &Pair{tiles: []tile.Tile{t, t}}, nil
}

func (p *Pair) ToTiles() []tile.Tile {
	return p.tiles
}
