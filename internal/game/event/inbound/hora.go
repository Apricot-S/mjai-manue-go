package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Hora struct {
	Actor      int `validate:"min=0,max=3"`
	Target     int `validate:"min=0,max=3"`
	Pai        *base.Pai
	HoraPoints *int `validate:"omitnil,min=0"`
	Scores     *[4]int
}

func NewHora(actor int, target int, pai *base.Pai, horaPoints *int, scores *[4]int) (*Hora, error) {
	h := &Hora{
		Actor:      actor,
		Target:     target,
		Pai:        pai,
		HoraPoints: horaPoints,
		Scores:     scores,
	}

	if h.Pai != nil && h.Pai.IsUnknown() {
		return nil, fmt.Errorf("dora marker cannot be unknown: %s", h.Pai.ToString())
	}

	if err := eventValidator.Struct(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Hora) isInboundEvent() {}
