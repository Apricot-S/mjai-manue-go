package inbound

import (
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

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

	switch header.Type {
	case "start_kyoku":
		return ParseStartKyoku(b)
	case "tsumo":
		return ParseTsumo(b)
	case "dahai":
		return ParseDahai(b)
	default:
		return nil, fmt.Errorf("unsupported message type: %q", header.Type)
	}
}
