package tile

import (
	"fmt"
	"slices"
)

// NumTileType34 is the number of distinct tile types without red fives.
const NumTileType34 = 3*9 + 4 + 3

// NumTileType37 is the number of distinct tile types with red fives.
const NumTileType37 = NumTileType34 + 3

// NumTileType38 is the number of distinct tile types with red fives and unknown tile (?).
const NumTileType38 = NumTileType37 + 1

const minTileID = 0

const minSuitsID = minTileID
const maxSuitsID = minSuitsID + 9*3 - 1
const minHonorsID = maxSuitsID + 1
const maxHonorsID = minHonorsID + 4 + 3 - 1
const minRedID = maxHonorsID + 1
const maxRedID = minRedID + 2
const unknownID = maxRedID + 1

var tileCodes = [NumTileType38]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", // m
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p", // p
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s", // s
	"E", "S", "W", "N", "P", "F", "C", // z
	"5mr", "5pr", "5sr", // red
	"?", // unknown
}

var tileNumbers = [NumTileType38]int{
	1, 2, 3, 4, 5, 6, 7, 8, 9,
	1, 2, 3, 4, 5, 6, 7, 8, 9,
	1, 2, 3, 4, 5, 6, 7, 8, 9,
	1, 2, 3, 4, 5, 6, 7,
	5, 5, 5,
	0,
}

var tileIsReds = [NumTileType38]bool{
	false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false, false, false,
	false, false, false, false, false, false, false,
	true, true, true,
	false,
}

type Tile struct {
	id     int
	number int
	isRed  bool
}

func (t *Tile) ID() int {
	return t.id
}

func NewTileFromID(id int) (*Tile, error) {
	if id < minTileID || id >= NumTileType38 {
		return nil, fmt.Errorf("invalid tile id: %d", id)
	}
	return &Tile{
		id:     id,
		number: tileNumbers[id],
		isRed:  tileIsReds[id],
	}, nil
}

func MustTileFromID(id int) *Tile {
	if id < minTileID || id >= NumTileType38 {
		panic(fmt.Sprintf("invalid tile id: %d", id))
	}
	return &Tile{
		id:     id,
		number: tileNumbers[id],
		isRed:  tileIsReds[id],
	}
}

func NewTileFromCode(code string) (*Tile, error) {
	id := slices.Index(tileCodes[:], code)
	if id == -1 {
		return nil, fmt.Errorf("invalid tile code: %s", code)
	}
	return NewTileFromID(id)
}

func MustTileFromCode(code string) *Tile {
	id := slices.Index(tileCodes[:], code)
	if id == -1 {
		panic(fmt.Sprintf("invalid tile code: %s", code))
	}
	return MustTileFromID(id)
}

func (t *Tile) Code() string {
	return tileCodes[t.id]
}

func (t *Tile) Number() int {
	return t.number
}

func (t *Tile) IsRed() bool {
	return t.isRed
}

func (t *Tile) IsSuits() bool {
	return t.id < minHonorsID || t.IsRed()
}

func (t *Tile) IsHonors() bool {
	return !t.IsSuits() && !t.IsUnknown()
}

func (t *Tile) IsUnknown() bool {
	return t.id == unknownID
}
