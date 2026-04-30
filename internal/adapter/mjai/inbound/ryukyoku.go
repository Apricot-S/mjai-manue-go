package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

type Ryukyoku struct {
	Type    string     `json:"type"`
	Reason  string     `json:"reason,omitempty"`
	Tehais  [][]string `json:"tehais,omitempty"`
	Tenpais []bool     `json:"tenpais,omitempty"`
	Deltas  []int      `json:"deltas,omitempty"`
	Scores  []int      `json:"scores,omitempty"`
}

func (*Ryukyoku) inboundMessage() {}

func (m *Ryukyoku) ToEvent() (*event.DrawRound, error) {
	if m == nil {
		return nil, fmt.Errorf("ryukyoku message is nil")
	}
	if m.Type != "ryukyoku" {
		return nil, fmt.Errorf("unexpected message type: %q", m.Type)
	}
	tenpais, err := parseOptionalTenpaisField(m.Tenpais)
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
	return event.NewDrawRound(m.Reason, tenpais, deltas, scores), nil
}
