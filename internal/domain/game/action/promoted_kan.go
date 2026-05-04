package action

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type PromotedKan struct {
	actor    seat.Seat
	added    tile.Tile
	consumed [3]tile.Tile
}

func NewPromotedKan(actor seat.Seat, added tile.Tile, consumed [3]tile.Tile) (*PromotedKan, error) {
	if added.IsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for PromotedKan action added tile: %s", added)
	}
	if tile.Tiles(consumed[:]).ContainsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for PromotedKan action consumed tiles: %v", consumed)
	}
	return &PromotedKan{actor: actor, added: added, consumed: consumed}, nil
}

func (*PromotedKan) isAction() {}

func (k *PromotedKan) Actor() seat.Seat {
	return k.actor
}

func (k *PromotedKan) Added() tile.Tile {
	return k.added
}

func (k *PromotedKan) Consumed() [3]tile.Tile {
	return k.consumed
}
