package block

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Block interface {
	ToTiles() []tile.Tile
}

type Sequence struct {
	tiles [3]tile.Tile
}

func NewSequence(t tile.Tile) (*Sequence, error) {
	if !t.IsSuits() {
		return nil, fmt.Errorf("cannot create sequence from honors or unknown tile: %s", t.Code())
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create sequence from red five")
	}
	if t.Number() > 7 {
		return nil, fmt.Errorf("cannot create sequence starting with 8 or 9: %s", t.Code())
	}
	return &Sequence{tiles: [3]tile.Tile{t, *t.Next(1), *t.Next(2)}}, nil
}

func MustSequence(t tile.Tile) *Sequence {
	s, err := NewSequence(t)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Sequence) ToTiles() []tile.Tile {
	return s.tiles[:]
}
