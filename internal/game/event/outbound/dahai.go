package outbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Dahai struct {
	action
	Pai       base.Pai
	Tsumogiri bool
}

func NewDahai(actor int, pai base.Pai, tsumogiri bool, log string) (*Dahai, error) {
	event := &Dahai{
		action: action{
			Actor: actor,
			Log:   log,
		},
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
