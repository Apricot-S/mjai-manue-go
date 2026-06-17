package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Pon struct {
	actor    seat.Seat
	target   seat.Seat
	taken    tile.Tile
	consumed [2]tile.Tile
}

func NewPon(actor, target seat.Seat, taken tile.Tile, consumed [2]tile.Tile) *Pon {
	return &Pon{actor: actor, target: target, taken: taken, consumed: consumed}
}

func (*Pon) isEvent() {}

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
