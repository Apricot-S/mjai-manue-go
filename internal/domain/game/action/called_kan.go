package action

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type CalledKan struct {
	actor    seat.Seat
	target   seat.Seat
	taken    tile.Tile
	consumed [3]tile.Tile
}

func NewCalledKan(actor, target seat.Seat, taken tile.Tile, consumed [3]tile.Tile) (*CalledKan, error) {
	if taken.IsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for CalledKan action taken tile: %s", taken)
	}
	if tile.Tiles(consumed[:]).ContainsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for CalledKan action consumed tiles: %v", consumed)
	}
	return &CalledKan{actor: actor, target: target, taken: taken, consumed: consumed}, nil
}

func (*CalledKan) isAction() {}

func (k *CalledKan) Actor() seat.Seat {
	return k.actor
}

func (k *CalledKan) Target() seat.Seat {
	return k.target
}

func (k *CalledKan) Taken() tile.Tile {
	return k.taken
}

func (k *CalledKan) Consumed() [3]tile.Tile {
	return k.consumed
}
