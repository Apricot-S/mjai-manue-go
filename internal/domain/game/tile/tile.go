package tile

import (
	"fmt"
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
const minWindID = minHonorsID
const maxWindID = minWindID + 4 - 1
const minDragonID = maxWindID + 1
const maxDragonID = minDragonID + 3 - 1
const maxHonorsID = maxDragonID
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

var tileCodeToID = func() map[string]int {
	m := make(map[string]int, NumTileType38)
	for id, code := range tileCodes {
		m[code] = id
	}
	return m
}()

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

var tileIsYaochus = func() [NumTileType38]bool {
	ys := [NumTileType38]bool{}
	for _, id := range YaochuhaiIDs {
		ys[id] = true
	}
	return ys
}()

var doraIndicatorToDora = [NumTileType38]int{
	1, 2, 3, 4, 5, 6, 7, 8, 0,
	10, 11, 12, 13, 14, 15, 16, 17, 9,
	19, 20, 21, 22, 23, 24, 25, 26, 18,
	28, 29, 30, 27, 32, 33, 31,
	5, 13, 22,
	37,
}

type Tile struct {
	id int8
}

var tilesByID = func() [NumTileType38]Tile {
	ts := [NumTileType38]Tile{}
	for id := range ts {
		ts[id] = Tile{
			id: int8(id),
		}
	}
	return ts
}()

func newTileFromValidID(id int) Tile {
	return tilesByID[id]
}

func MustTileFromID(id int) Tile {
	if id < minTileID || id >= NumTileType38 {
		panic(fmt.Sprintf("invalid tile id: %d", id))
	}
	return newTileFromValidID(id)
}

func NewTileFromCode(code string) (Tile, error) {
	id, ok := tileCodeToID[code]
	if !ok {
		return Tile{}, fmt.Errorf("invalid tile code: %s", code)
	}
	return newTileFromValidID(id), nil
}

func MustTileFromCode(code string) Tile {
	id, ok := tileCodeToID[code]
	if !ok {
		panic(fmt.Sprintf("invalid tile code: %s", code))
	}
	return newTileFromValidID(id)
}

func (t Tile) ID() int {
	return int(t.id)
}

func (t Tile) String() string {
	return tileCodes[t.ID()]
}

func (t Tile) Color() rune {
	return tileColors[t.ID()]
}

func (t Tile) Number() int {
	return tileNumbers[t.ID()]
}

func (t Tile) IsRed() bool {
	return tileIsReds[t.ID()]
}

func (t Tile) IsSuits() bool {
	return !t.IsHonors() && !t.IsUnknown()
}

func (t Tile) IsHonors() bool {
	return t.Color() == HonorsColor
}

func (t Tile) IsWind() bool {
	return minWindID <= t.ID() && t.ID() <= maxWindID
}

func (t Tile) IsDragon() bool {
	return minDragonID <= t.ID() && t.ID() <= maxDragonID
}

func (t Tile) IsYaochu() bool {
	return tileIsYaochus[t.ID()]
}

func (t Tile) IsUnknown() bool {
	return t.ID() == unknownID
}

func (t Tile) Next(n int) *Tile {
	if t.IsUnknown() || t.IsHonors() {
		return nil
	}

	nextNumber := t.Number() + n
	if nextNumber < 1 || 9 < nextNumber {
		return nil
	}

	nextID := t.RemoveRed().ID() + n
	return new(MustTileFromID(nextID))
}

func (t Tile) NextForDora() Tile {
	return MustTileFromID(doraIndicatorToDora[t.ID()])
}

func (t Tile) AddRed() Tile {
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

func (t Tile) RemoveRed() Tile {
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

func (t Tile) HasSameSymbol(other Tile) bool {
	return t.Number() == other.Number() && t.Color() == other.Color()
}
