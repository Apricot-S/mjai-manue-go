package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Kakan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

func (*Kakan) inboundMessage() {}

func (m *Kakan) ToEvent() (*event.PromotedKan, error) {
	if m == nil {
		return nil, fmt.Errorf("kakan message is nil")
	}
	if m.Type != "kakan" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}
	added, err := parseKnownTileField("pai", m.Pai)
	if err != nil {
		return nil, err
	}
	consumed, err := parseConsumed3(m.Consumed)
	if err != nil {
		return nil, err
	}
	return event.NewPromotedKan(*actor, *added, consumed), nil
}
