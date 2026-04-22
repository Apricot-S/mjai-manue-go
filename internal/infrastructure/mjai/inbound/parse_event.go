package inbound

import (
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

var parseEventByType = map[string]func([]byte) (event.Event, error){
	"start_kyoku": parseToEvent[*StartKyoku](),
	"tsumo":       parseToEvent[*Tsumo](),
	"dahai":       parseToEvent[*Dahai](),
}

// ParseEvent converts a single mjai inbound message (one JSON object) into a domain event.
//
// It dispatches by the "type" field. Unknown message types return an error.
//
// This function is intentionally pure with respect to I/O: it does not read from
// any stream and only depends on the given bytes. Framing (e.g. JSON Lines) is a
// transport concern.
func ParseEvent(b []byte) (event.Event, error) {
	var header struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(b, &header); err != nil {
		return nil, err
	}
	if header.Type == "" {
		return nil, fmt.Errorf("message type is missing")
	}

	parser, ok := parseEventByType[header.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported message type: %q", header.Type)
	}
	return parser(b)
}

type ToEventer[E event.Event] interface {
	ToEvent() (E, error)
}

func parseToEvent[M ToEventer[E], E event.Event]() func([]byte) (event.Event, error) {
	return func(b []byte) (event.Event, error) {
		var msg M
		if err := json.Unmarshal(b, &msg); err != nil {
			return nil, err
		}
		ev, err := msg.ToEvent()
		if err != nil {
			return nil, err
		}
		return ev, nil
	}
}
