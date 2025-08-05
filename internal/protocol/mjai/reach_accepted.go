package mjai

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type ReachAccepted struct {
	Message
	Actor  int   `json:"actor" validate:"min=0,max=3"`
	Scores []int `json:"scores,omitempty"`
}

func NewReachAccepted(actor int, scores []int) (*ReachAccepted, error) {
	if scores != nil && len(scores) != 4 {
		return nil, fmt.Errorf("invalid number of scores: %v", scores)
	}

	m := &ReachAccepted{
		Message: Message{Type: TypeReachAccepted},
		Actor:   actor,
		Scores:  scores,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *ReachAccepted) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeReachAccepted {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner ReachAccepted
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *ReachAccepted) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner ReachAccepted
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (ReachAccepted)(mm)
	if m.Type != TypeReachAccepted {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}

	return messageValidator.Struct(m)
}
