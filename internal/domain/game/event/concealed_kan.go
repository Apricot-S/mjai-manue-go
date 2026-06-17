package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type ConcealedKan struct {
	actor    seat.Seat
	consumed [4]tile.Tile
}

func NewConcealedKan(actor seat.Seat, consumed [4]tile.Tile) *ConcealedKan {
	return &ConcealedKan{actor: actor, consumed: consumed}
}

func (*ConcealedKan) isEvent() {}

func (k *ConcealedKan) Actor() seat.Seat {
	return k.actor
}

func (k *ConcealedKan) Consumed() [4]tile.Tile {
	return k.consumed
}
