package meld

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type PromotedKan struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	added    tile.Tile
	target   seat.Seat
	tiles    []tile.Tile
}

func NewPromotedKan(
	taken tile.Tile,
	consumed [2]tile.Tile,
	added tile.Tile,
	target seat.Seat,
) (*PromotedKan, error) {
	tiles := tile.Tiles{taken, consumed[0], consumed[1], added}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return t.IsUnknown() }) {
		return nil, fmt.Errorf("unknown tile cannot use for Promoted Kan")
	}
	if slices.ContainsFunc(tiles, func(t tile.Tile) bool { return !taken.HasSameSymbol(&t) }) {
		return nil, fmt.Errorf("mismatch taken: %+v, consumed: %+v, added: %+v", taken, consumed, added)
	}
	if taken.IsSuits() && taken.Number() == 5 && countRed(tiles) != 1 {
		return nil, fmt.Errorf("must contain a red five for Promoted Kan of 5; taken: %+v, consumed: %+v, added: %+v", taken, consumed, added)
	}

	tiles.Sort()

	csm := tile.Tiles(consumed[:])
	csm.Sort()

	return &PromotedKan{
		taken:    taken,
		consumed: [2]tile.Tile(csm),
		added:    added,
		target:   target,
		tiles:    tiles,
	}, nil
}

func MustPromotedKan(
	taken tile.Tile,
	consumed [2]tile.Tile,
	added tile.Tile,
	target seat.Seat,
) *PromotedKan {
	k, err := NewPromotedKan(taken, consumed, added, target)
	if err != nil {
		panic(err)
	}
	return k
}

func (k *PromotedKan) Taken() *tile.Tile {
	return &k.taken
}

func (k *PromotedKan) Consumed() []tile.Tile {
	return k.consumed[:]
}

func (k *PromotedKan) Added() *tile.Tile {
	return &k.added
}

func (k *PromotedKan) Target() *seat.Seat {
	return &k.target
}

func (k *PromotedKan) ToTiles() []tile.Tile {
	return k.tiles
}

func (k PromotedKan) String() string {
	taken := k.Taken().String()
	target := k.Target().Index()
	consumed0 := k.consumed[0].String()
	consumed1 := k.consumed[1].String()
	added := k.Added().String()

	return fmt.Sprintf("[%s(%d)/%s %s %s]", taken, target, consumed0, consumed1, added)
}
