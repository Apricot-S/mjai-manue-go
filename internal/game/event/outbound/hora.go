package outbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Hora struct {
	action
	Target int `validate:"min=0,max=3"`
	Pai    base.Pai
}

func NewHora(actor int, target int, pai base.Pai, log string) (*Hora, error) {
	event := &Hora{
		action: action{
			Actor: actor,
			Log:   log,
		},
		Target: target,
		Pai:    pai,
	}

	if event.Pai.IsUnknown() {
		return nil, fmt.Errorf("hora tile must not be unknown: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
