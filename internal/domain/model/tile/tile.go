package tile

import (
	"fmt"
	"slices"
)

const NumTileType = 3*9 + 4 + 3 + 3 + 1

var tileCodes = [NumTileType]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", // m
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p", // p
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s", // s
	"E", "S", "W", "N", "P", "F", "C", // z
	"5mr", "5pr", "5sr", // red
	"?", // unknown
}

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

func NewTileFromCode(code string) (*Tile, error) {
	id := slices.Index(tileCodes[:], code)
	if id == -1 {
		return nil, fmt.Errorf("invalid tile code: %s", code)
	}
	return NewTileFromID(id)
}

func (t *Tile) Code() string {
	return tileCodes[t.id]
}
