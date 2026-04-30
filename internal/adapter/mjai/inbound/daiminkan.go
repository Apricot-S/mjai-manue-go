package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Daiminkan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

func (*Daiminkan) inboundMessage() {}

func (m *Daiminkan) ToEvent() (*event.CalledKan, error) {
	if m == nil {
		return nil, fmt.Errorf("daiminkan message is nil")
	}
	if m.Type != "daiminkan" {
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
	taken, err := parseKnownTileField("pai", m.Pai)
	if err != nil {
		return nil, err
	}
	consumed, err := parseConsumed3(m.Consumed)
	if err != nil {
		return nil, err
	}
	return event.NewCalledKan(*actor, *target, *taken, consumed), nil
}
