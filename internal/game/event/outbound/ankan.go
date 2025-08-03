package outbound

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Ankan struct {
	action
	Consumed [4]base.Pai
}

func NewAnkan(actor int, consumed [4]base.Pai, log string) (*Ankan, error) {
	event := &Ankan{
		action: action{
			Actor: actor,
			Log:   log,
		},
		Consumed: consumed,
	}

	isValidPais := !slices.ContainsFunc(event.Consumed[1:], func(p base.Pai) bool {
		return !event.Consumed[0].HasSameSymbol(&p)
	})
	if !isValidPais {
		return nil, fmt.Errorf("all consumed tiles must be the same tile: %v", event)
	}

	isUnknown := slices.ContainsFunc(event.Consumed[:], func(p base.Pai) bool {
		return p.IsUnknown()
	})
	if isUnknown {
		return nil, fmt.Errorf("ankan tiles must not be unknown: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
