package outbound

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Kakan struct {
	action
	// Target for the original Pon
	Target int `validate:"min=0,max=3"`
	// Taken tile for the original Pon
	Taken    base.Pai
	Consumed [2]base.Pai
	// Tile added from hand to form Kakan
	Added base.Pai
}

func NewKakan(actor int, target int, taken base.Pai, consumed [2]base.Pai, added base.Pai, log string) (*Kakan, error) {
	event := &Kakan{
		action: action{
			Actor: actor,
			Log:   log,
		},
		Target:   target,
		Taken:    taken,
		Consumed: consumed,
		Added:    added,
	}

	if event.Actor == event.Target {
		return nil, fmt.Errorf("actor and target cannot be the same: %d", event.Actor)
	}

	var pais base.Pais = append(event.Consumed[:], event.Taken)
	isValidPais := !slices.ContainsFunc(pais, func(p base.Pai) bool {
		return !event.Added.HasSameSymbol(&p)
	})
	if !isValidPais {
		return nil, fmt.Errorf("all tiles must be the same tile: %v", event)
	}

	if event.Added.IsUnknown() {
		return nil, fmt.Errorf("kakan tiles must not be unknown: %v", event)
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
