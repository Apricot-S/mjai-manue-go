package model

import "fmt"

type Tile struct {
	id int
}

func (t *Tile) ID() int {
	return t.id
}

func NewTileFromID(id int) (*Tile, error) {
	if id < 0 || id >= 38 {
		return nil, fmt.Errorf("invalid tile id: %d", id)
	}
	return &Tile{id: id}, nil
}
