package action

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"

type Pass struct {
	actor seat.Seat
}

func NewPass(actor seat.Seat) *Pass {
	return &Pass{actor: actor}
}

func (*Pass) isAction() {}

func (p *Pass) Actor() seat.Seat {
	return p.actor
}
