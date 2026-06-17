package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Reach struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
}

func (*Reach) inboundMessage() {}

func (m *Reach) ToEvent() (*event.Riichi, error) {
	if m == nil {
		return nil, fmt.Errorf("reach message is nil")
	}
	if m.Type != "reach" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}
	return event.NewRiichi(*actor), nil
}
