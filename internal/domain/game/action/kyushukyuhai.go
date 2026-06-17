package action

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"

type Kyushukyuhai struct {
	actor seat.Seat
}

func NewKyushukyuhai(actor seat.Seat) *Kyushukyuhai {
	return &Kyushukyuhai{actor: actor}
}

func (*Kyushukyuhai) isAction() {}

func (k *Kyushukyuhai) Actor() seat.Seat {
	return k.actor
}
