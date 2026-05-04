package action

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Chii struct {
	actor    seat.Seat
	target   seat.Seat
	taken    tile.Tile
	consumed [2]tile.Tile
}

func NewChii(actor, target seat.Seat, taken tile.Tile, consumed [2]tile.Tile) (*Chii, error) {
	if taken.IsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for Chii action taken tile: %s", taken)
	}
	if tile.Tiles(consumed[:]).ContainsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for Chii action consumed tiles: %v", consumed)
	}
	return &Chii{actor: actor, target: target, taken: taken, consumed: consumed}, nil
}

func (*Chii) isAction() {}

func (c *Chii) Actor() seat.Seat {
	return c.actor
}

func (c *Chii) Target() seat.Seat {
	return c.target
}

func (c *Chii) Taken() tile.Tile {
	return c.taken
}

func (c *Chii) Consumed() [2]tile.Tile {
	return c.consumed
}
