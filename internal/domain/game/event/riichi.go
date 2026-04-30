package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type Riichi struct {
	actor seat.Seat
}

func NewRiichi(actor seat.Seat) *Riichi {
	return &Riichi{actor: actor}
}

func (*Riichi) isEvent() {}

func (r *Riichi) Actor() seat.Seat {
	return r.actor
}
