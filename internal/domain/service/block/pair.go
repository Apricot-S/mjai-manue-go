package block

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

type Pair struct {
	tiles [2]tile.Tile
}

func NewPair(t tile.Tile) (*Pair, error) {
	if t.IsUnknown() {
		return nil, fmt.Errorf("cannot create pair from unknown tile")
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create pair from red five")
	}
	return &Pair{tiles: [2]tile.Tile{t, t}}, nil
}

func MustPair(t tile.Tile) *Pair {
	p, err := NewPair(t)
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Pair) ToTiles() []tile.Tile {
	return p.tiles[:]
}
