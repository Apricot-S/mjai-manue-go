package inbound

import (
	"encoding/json/v2"
	"fmt"
	"io"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

// ParseEvent converts a single mjai inbound message (one JSON object) into a domain event.
//
// It dispatches by the "type" field. Unknown message types return an error.
func ParseEvent(r io.Reader) (event.Event, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

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
		return parseStartKyokuBytes(b)
	case "tsumo":
		return parseTsumoBytes(b)
	default:
		return nil, fmt.Errorf("unsupported message type: %q", header.Type)
	}
}
