package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Dora struct {
	Type       string `json:"type"`
	DoraMarker string `json:"dora_marker"`
}

func (*Dora) inboundMessage() {}

func (m *Dora) ToEvent() (*event.Dora, error) {
	if m == nil {
		return nil, fmt.Errorf("dora message is nil")
	}
	if m.Type != "dora" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	indicator, err := parseTileField("dora_marker", m.DoraMarker)
	if err != nil {
		return nil, err
	}
	return event.NewDora(*indicator), nil
}
