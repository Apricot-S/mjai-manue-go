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

type Triplet struct {
	tiles []tile.Tile
}

func NewTriplet(t tile.Tile) (*Triplet, error) {
	if t.IsUnknown() {
		return nil, fmt.Errorf("cannot create triplet from unknown tile")
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create triplet from red five")
	}
	return &Triplet{tiles: []tile.Tile{t, t, t}}, nil
}

func (p *Triplet) ToTiles() []tile.Tile {
	return p.tiles
}

type Quad struct {
	tiles []tile.Tile
}

func NewQuad(t tile.Tile) (*Quad, error) {
	if t.IsUnknown() {
		return nil, fmt.Errorf("cannot create quad from unknown tile")
	}
	if t.IsRed() {
		return nil, fmt.Errorf("cannot create quad from red five")
	}
	return &Quad{tiles: []tile.Tile{t, t, t, t}}, nil
}

func (p *Quad) ToTiles() []tile.Tile {
	return p.tiles
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
