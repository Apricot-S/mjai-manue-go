package inbound

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Ankan struct {
	Actor    int `validate:"min=0,max=3"`
	Consumed [4]base.Pai
}

func NewAnkan(actor int, consumed [4]base.Pai) (*Ankan, error) {
	event := &Ankan{
		Actor:    actor,
		Consumed: consumed,
	}

	isValidPais := !slices.ContainsFunc(event.Consumed[1:], func(p base.Pai) bool {
		return !event.Consumed[0].HasSameSymbol(&p)
	})
	if !isValidPais {
		return nil, fmt.Errorf("all consumed tiles must be the same tile: %v", event)
	}

	if event.Consumed[0].IsUnknown() {
		return nil, fmt.Errorf("ankan tiles must not be unknown: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
