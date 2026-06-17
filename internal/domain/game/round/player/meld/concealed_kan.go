package meld

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type ConcealedKan struct {
	consumed [4]tile.Tile
}

func NewConcealedKan(consumed [4]tile.Tile) (*ConcealedKan, error) {
	var csm tile.Tiles = slices.Clone(consumed[:])
	c0 := &csm[0]

	if csm.ContainsUnknown() {
		return nil, fmt.Errorf("unknown tile cannot use for Concealed Kan")
	}
	if csm.CountSameSymbol(*c0) != 4 {
		return nil, fmt.Errorf("mismatch consumed: %+v", consumed)
	}
	if c0.IsSuits() && c0.Number() == 5 && countRed(csm) != 1 {
		return nil, fmt.Errorf("must contain a red five for Concealed Kan of 5; consumed: %+v", consumed)
	}

	csm.Sort()

	return &ConcealedKan{consumed: [4]tile.Tile(csm)}, nil
}

func MustConcealedKan(consumed [4]tile.Tile) *ConcealedKan {
	k, err := NewConcealedKan(consumed)
	if err != nil {
		panic(err)
	}
	return k
}

func (k *ConcealedKan) Consumed() []tile.Tile {
	return k.consumed[:]
}

func (k *ConcealedKan) ToTiles() []tile.Tile {
	return k.consumed[:]
}

func (k ConcealedKan) String() string {
	// Red five is in consumed[3]
	return fmt.Sprintf("[# %s %s #]", k.consumed[2], k.consumed[3])
}
