package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

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

	if event.Pai.IsUnknown() {
		return nil, fmt.Errorf("dahai must not be unknown: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (d *Dahai) isInboundEvent() {}
