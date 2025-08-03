package inbound

import "github.com/Apricot-S/mjai-manue-go/internal/base"

type Dahai struct {
	Actor     int `validate:"min=0,max=3"`
	Pai       base.Pai
	Tsumogiri bool
}

func NewDahai(actor int, pai base.Pai, tsumogiri bool) (*Dahai, error) {
	event := &Dahai{
		Actor:     actor,
		Pai:       pai,
		Tsumogiri: tsumogiri,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Dahai) isInboundEvent() {}
