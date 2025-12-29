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

const ManzuColor = 'm'
const PinzuColor = 'p'
const SouzuColor = 's'
const HonorsColor = 't'
const UnknownColor = '?'

var YaochuhaiIDs = [13]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

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

var tileColors = [NumTileType38]rune{
	'm', 'm', 'm', 'm', 'm', 'm', 'm', 'm', 'm',
	'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
	's', 's', 's', 's', 's', 's', 's', 's', 's',
	't', 't', 't', 't', 't', 't', 't',
	'm', 'p', 's',
	'?',
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
	color  rune
	number int
	isRed  bool
}

func NewTileFromID(id int) (*Tile, error) {
	if id < minTileID || id >= NumTileType38 {
		return nil, fmt.Errorf("invalid tile id: %d", id)
	}
	return &Tile{
		id:     id,
		color:  tileColors[id],
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
		color:  tileColors[id],
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

func (t *Tile) ID() int {
	return t.id
}

func (t Tile) String() string {
	return tileCodes[t.id]
}

func (t *Tile) Color() rune {
	return t.color
}

func (t *Tile) Number() int {
	return t.number
}

func (t *Tile) IsRed() bool {
	return t.isRed
}

func (t *Tile) IsSuits() bool {
	return !t.IsHonors() && !t.IsUnknown()
}

func (t *Tile) IsHonors() bool {
	return t.Color() == HonorsColor
}

func (t *Tile) IsYaochu() bool {
	return slices.Contains(YaochuhaiIDs[:], t.ID())
}

func (t *Tile) IsUnknown() bool {
	return t.id == unknownID
}

func (t *Tile) Next(n int) *Tile {
	if t.IsUnknown() || t.IsHonors() {
		return nil
	}

	nextNumber := t.Number() + n
	if nextNumber < 1 || 9 < nextNumber {
		return nil
	}

	nextID := t.RemoveRed().ID() + n
	return MustTileFromID(nextID)
}

func (t *Tile) NextForDora() *Tile {
	if t.IsUnknown() {
		return t
	}

	panic("unimplemented!")
}

func (t *Tile) AddRed() *Tile {
	switch t.ID() {
	case 4:
		return MustTileFromID(minRedID)
	case 13:
		return MustTileFromID(minRedID + 1)
	case 22:
		return MustTileFromID(minRedID + 2)
	default:
		return t
	}
}

func (t *Tile) RemoveRed() *Tile {
	switch t.ID() {
	case minRedID:
		return MustTileFromID(4)
	case minRedID + 1:
		return MustTileFromID(13)
	case minRedID + 2:
		return MustTileFromID(22)
	default:
		return t
	}
}

func (t *Tile) HasSameSymbol(other *Tile) bool {
	return t.Number() == other.Number() && t.Color() == other.Color()
}
