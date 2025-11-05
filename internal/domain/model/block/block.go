package block

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Block interface {
	ToTiles() []tile.Tile
}

type Sequence struct {
	tiles []tile.Tile
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
	return &Sequence{tiles: []tile.Tile{t, *t.Next(1), *t.Next(2)}}, nil
}

func (p *Sequence) ToTiles() []tile.Tile {
	return p.tiles
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
