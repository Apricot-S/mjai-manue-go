package inbound

import "github.com/Apricot-S/mjai-manue-go/internal/base"

type Tsumo struct {
	Actor int `validate:"min=0,max=3"`
	Pai   base.Pai
}

func NewTsumo(actor int, pai base.Pai) (*Tsumo, error) {
	event := &Tsumo{
		Actor: actor,
		Pai:   pai,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (t *Tsumo) isInboundEvent() {}
