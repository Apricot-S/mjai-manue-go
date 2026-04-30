package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type PromotedKan struct {
	actor    seat.Seat
	added    tile.Tile
	consumed [3]tile.Tile
}

func NewPromotedKan(actor seat.Seat, added tile.Tile, consumed [3]tile.Tile) *PromotedKan {
	return &PromotedKan{actor: actor, added: added, consumed: consumed}
}

func (*PromotedKan) isEvent() {}

func (k *PromotedKan) Actor() seat.Seat {
	return k.actor
}

func (k *PromotedKan) Added() tile.Tile {
	return k.added
}

func (k *PromotedKan) Consumed() [3]tile.Tile {
	return k.consumed
}
