package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type EndKyoku struct {
	Type string `json:"type"`
}

func (*EndKyoku) inboundMessage() {}

func (m *EndKyoku) ToEvent() (*event.EndRound, error) {
	if m == nil {
		return nil, fmt.Errorf("end kyoku message is nil")
	}
	if m.Type != "end_kyoku" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	return event.NewEndRound(), nil
}
