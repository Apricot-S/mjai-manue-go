package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

// ParseEvent converts a decoded mjai inbound message into a domain event.
//
// Not all inbound messages correspond to domain events. Those messages return an error.
func ParseEvent(msg Message) (event.Event, error) {
	switch m := msg.(type) {
	case *StartKyoku:
		return m.ToEvent()
	case *Tsumo:
		return m.ToEvent()
	case *Dahai:
		return m.ToEvent()
	default:
		return nil, fmt.Errorf("message cannot be converted to event: %T", msg)
	}
}
