package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

// Since mjai-manue does not select Nine Different Terminals and Honors (九種九牌),
// it does not need to be an Action.
type Ryukyoku struct {
	Message
	Scores []int `json:"scores,omitempty"`
}

func NewRyukyoku(scores []int) (*Ryukyoku, error) {
	if scores != nil && len(scores) != 4 {
		return nil, fmt.Errorf("invalid number of scores: %v", scores)
	}

	m := &Ryukyoku{
		Message: Message{Type: TypeRyukyoku},
		Scores:  scores,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Ryukyoku) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeRyukyoku {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Ryukyoku
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Ryukyoku) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Ryukyoku
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Ryukyoku)(mm)
	if m.Type != TypeRyukyoku {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}

	return messageValidator.Struct(m)
}

func (m *Ryukyoku) ToEvent() *inbound.Ryukyoku {
	var scores *[4]int = nil
	if m.Scores != nil {
		scores = (*[4]int)(m.Scores)
	}

	return inbound.NewRyukyoku(scores)
}
