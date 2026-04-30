package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Tsumo struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
	Pai   string `json:"pai"`
}

func (*Tsumo) inboundMessage() {}

func (m *Tsumo) ToEvent() (*event.Draw, error) {
	if m == nil {
		return nil, fmt.Errorf("tsumo message is nil")
	}
	if m.Type != "tsumo" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}

	pai, err := parseTileField("pai", m.Pai)
	if err != nil {
		return nil, err
	}

	return event.NewDraw(*actor, *pai), nil
}
