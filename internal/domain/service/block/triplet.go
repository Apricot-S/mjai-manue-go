package block

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

type Triplet struct {
	tiles [3]tile.Tile
}

func NewTriplet(t tile.Tile) (*Triplet, error) {
	if t.IsUnknown() {
		return nil, fmt.Errorf("cannot create triplet from unknown tile")
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create triplet from red five")
	}
	return &Triplet{tiles: [3]tile.Tile{t, t, t}}, nil
}

func MustTriplet(t tile.Tile) *Triplet {
	tr, err := NewTriplet(t)
	if err != nil {
		panic(err)
	}
	return tr
}

func (tr *Triplet) ToTiles() []tile.Tile {
	return tr.tiles[:]
}
