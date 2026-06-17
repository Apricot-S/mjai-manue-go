package action

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Pon struct {
	actor    seat.Seat
	target   seat.Seat
	taken    tile.Tile
	consumed [2]tile.Tile
}

func NewPon(actor, target seat.Seat, taken tile.Tile, consumed [2]tile.Tile) (*Pon, error) {
	if taken.IsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for Pon action taken tile: %s", taken)
	}
	if tile.Tiles(consumed[:]).ContainsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for Pon action consumed tiles: %v", consumed)
	}
	return &Pon{actor: actor, target: target, taken: taken, consumed: consumed}, nil
}

func (*Pon) isAction() {}

func (p *Pon) Actor() seat.Seat {
	return p.actor
}

func (p *Pon) Target() seat.Seat {
	return p.target
}

func (p *Pon) Taken() tile.Tile {
	return p.taken
}

func (p *Pon) Consumed() [2]tile.Tile {
	return p.consumed
}
