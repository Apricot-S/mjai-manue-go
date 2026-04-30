package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Dahai struct {
	Type      string `json:"type"`
	Actor     int    `json:"actor"`
	Pai       string `json:"pai"`
	Tsumogiri bool   `json:"tsumogiri"`
}

func (*Dahai) inboundMessage() {}

func (m *Dahai) ToEvent() (*event.Discard, error) {
	if m == nil {
		return nil, fmt.Errorf("dahai message is nil")
	}
	if m.Type != "dahai" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}

	pai, err := parseKnownTileField("pai", m.Pai)
	if err != nil {
		return nil, err
	}

	return event.NewDiscard(*actor, *pai, m.Tsumogiri), nil
}
