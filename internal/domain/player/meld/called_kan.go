package meld

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

type CalledKan struct {
	taken    tile.Tile
	consumed [3]tile.Tile
	target   playerid.PlayerID
	tiles    []tile.Tile
}

func NewCalledKan(taken tile.Tile, consumed [3]tile.Tile, target playerid.PlayerID) (*CalledKan, error) {
	tiles := tile.Tiles{taken, consumed[0], consumed[1], consumed[2]}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return t.IsUnknown() }) {
		return nil, fmt.Errorf("unknown tile cannot use for Called Kan")
	}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return !taken.HasSameSymbol(&t) }) {
		return nil, fmt.Errorf("mismatch taken: %+v, consumed: %+v", taken, consumed)
	}
	if taken.IsSuits() && taken.Number() == 5 && countRed(tiles) != 1 {
		return nil, fmt.Errorf("must contain a red five for Called Kan of 5; taken: %+v, consumed: %+v", taken, consumed)
	}

	sort.Sort(tiles)

	csm := tile.Tiles(consumed[:])
	sort.Sort(csm)

	return &CalledKan{
		taken:    taken,
		consumed: [3]tile.Tile(csm),
		target:   target,
		tiles:    tiles,
	}, nil
}

func MustCalledKan(taken tile.Tile, consumed [3]tile.Tile, target playerid.PlayerID) *CalledKan {
	k, err := NewCalledKan(taken, consumed, target)
	if err != nil {
		panic(err)
	}
	return k
}

func (k *CalledKan) Taken() *tile.Tile {
	return &k.taken
}

func (k *CalledKan) Consumed() []tile.Tile {
	return k.consumed[:]
}

func (k *CalledKan) Target() *playerid.PlayerID {
	return &k.target
}

func (k *CalledKan) ToTiles() []tile.Tile {
	return k.tiles
}

func (k CalledKan) String() string {
	return meldToString(&k)
}
