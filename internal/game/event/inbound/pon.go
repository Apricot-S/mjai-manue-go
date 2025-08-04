package inbound

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Pon struct {
	Actor    int `validate:"min=0,max=3"`
	Target   int `validate:"min=0,max=3"`
	Taken    base.Pai
	Consumed [2]base.Pai
}

func NewPon(actor int, target int, taken base.Pai, consumed [2]base.Pai) (*Pon, error) {
	event := &Pon{
		Actor:    actor,
		Target:   target,
		Taken:    taken,
		Consumed: consumed,
	}

	if event.Actor == event.Target {
		return nil, fmt.Errorf("actor and target cannot be the same: %d", event.Actor)
	}

	isValidPais := !slices.ContainsFunc(event.Consumed[:], func(p base.Pai) bool {
		return !event.Taken.HasSameSymbol(&p)
	})
	if !isValidPais {
		return nil, fmt.Errorf("taken tile must be the same as the consumed tile: %v", event)
	}

	if event.Taken.IsUnknown() {
		return nil, fmt.Errorf("pon tiles must not be unknown: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
