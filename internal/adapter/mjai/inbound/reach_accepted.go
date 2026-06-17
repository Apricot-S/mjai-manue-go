package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type ReachAccepted struct {
	Type   string `json:"type"`
	Actor  int    `json:"actor"`
	Deltas []int  `json:"deltas,omitempty"`
	Scores []int  `json:"scores,omitempty"`
}

func (*ReachAccepted) inboundMessage() {}

func (m *ReachAccepted) ToEvent() (*event.RiichiAccepted, error) {
	if m == nil {
		return nil, fmt.Errorf("reach_accepted message is nil")
	}
	if m.Type != "reach_accepted" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}

	actor, err := parseSeatField("actor", m.Actor)
	if err != nil {
		return nil, err
	}
	deltas, err := parseOptionalScoresField("deltas", m.Deltas)
	if err != nil {
		return nil, err
	}
	scores, err := parseOptionalScoresField("scores", m.Scores)
	if err != nil {
		return nil, err
	}
	return event.NewRiichiAccepted(*actor, deltas, scores), nil
}
