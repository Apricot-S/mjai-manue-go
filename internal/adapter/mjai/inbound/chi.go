package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Chi struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

func (*Chi) inboundMessage() {}

func (m *Chi) ToEvent() (*event.Chii, error) {
	if m == nil {
		return nil, fmt.Errorf("chi message is nil")
	}
	if m.Type != "chi" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}
	target, err := parseSeatField("target", m.Target)
	if err != nil {
		return nil, err
	}
	taken, err := parseTileField("pai", m.Pai)
	if err != nil {
		return nil, err
	}
	consumed, err := parseConsumed2(m.Consumed)
	if err != nil {
		return nil, err
	}
	return event.NewChii(*actor, *target, *taken, consumed), nil
}
