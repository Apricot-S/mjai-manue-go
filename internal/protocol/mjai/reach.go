package mjai

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Reach struct {
	Action
}

func NewReach(actor int, log string) (*Reach, error) {
	m := &Reach{
		Action: Action{
			Message: Message{Type: TypeReach},
			Actor:   actor,
			Log:     log,
		},
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Reach) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeReach {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Reach
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Reach) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Reach
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Reach)(mm)
	if m.Type != TypeReach {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
