package action

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type ConcealedKan struct {
	actor    seat.Seat
	consumed [4]tile.Tile
}

func NewConcealedKan(actor seat.Seat, consumed [4]tile.Tile) (*ConcealedKan, error) {
	if hasUnknownTile(consumed[:]) {
		return nil, fmt.Errorf("unknown tile is not allowed for ConcealedKan action consumed tiles: %v", consumed)
	}
	return &ConcealedKan{actor: actor, consumed: consumed}, nil
}

func (*ConcealedKan) isAction() {}

func (k *ConcealedKan) Actor() seat.Seat {
	return k.actor
}

func (k *ConcealedKan) Consumed() [4]tile.Tile {
	return k.consumed
}
