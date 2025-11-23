package meld

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/playerid"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type PromotedKan struct {
	taken    tile.Tile
	consumed [2]tile.Tile
	added    tile.Tile
	target   playerid.PlayerID
	tiles    []tile.Tile
}

func NewPromotedKan(
	taken tile.Tile,
	consumed [2]tile.Tile,
	added tile.Tile,
	target playerid.PlayerID,
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

	sort.Sort(tiles)

	csm := tile.Tiles(consumed[:])
	sort.Sort(csm)

	return &PromotedKan{
		taken:    taken,
		consumed: [2]tile.Tile(csm),
		added:    added,
		target:   target,
		tiles:    tiles,
	}, nil
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

func (k *PromotedKan) Target() *playerid.PlayerID {
	return &k.target
}

func (k *PromotedKan) ToTiles() []tile.Tile {
	return k.tiles
}

func (k *PromotedKan) ToBlock() block.Block {
	// Red five is sorted after normal, so RemoveRed() is not necessary.
	return block.MustQuad(k.tiles[0])
}

func (k *PromotedKan) ToString() string {
	taken := k.Taken().Code()
	target := k.Target().Index()
	consumed0 := k.consumed[0].Code()
	consumed1 := k.consumed[1].Code()
	added := k.Added().Code()

	return fmt.Sprintf("[%s(%d)/%s %s %s]", taken, target, consumed0, consumed1, added)
}
