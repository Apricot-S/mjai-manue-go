package block

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Quad struct {
	tiles [4]tile.Tile
}

func NewQuad(t tile.Tile) (*Quad, error) {
	if t.IsUnknown() {
		return nil, fmt.Errorf("cannot create quad from unknown tile")
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create quad from red five")
	}
	return &Quad{tiles: [4]tile.Tile{t, t, t, t}}, nil
}

func MustQuad(t tile.Tile) *Quad {
	q, err := NewQuad(t)
	if err != nil {
		panic(err)
	}
	return q
}

func (q *Quad) ToTiles() []tile.Tile {
	return q.tiles[:]
}
