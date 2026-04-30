package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Ankan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Consumed []string `json:"consumed"`
}

func (*Ankan) inboundMessage() {}

func (m *Ankan) ToEvent() (*event.ConcealedKan, error) {
	if m == nil {
		return nil, fmt.Errorf("ankan message is nil")
	}
	if m.Type != "ankan" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}
	consumed, err := parseConsumed4(m.Consumed)
	if err != nil {
		return nil, err
	}
	return event.NewConcealedKan(*actor, consumed), nil
}
