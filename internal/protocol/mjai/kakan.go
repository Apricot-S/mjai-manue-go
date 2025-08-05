package mjai

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Kakan struct {
	Action
	Pai      string    `json:"pai" validate:"tile"`
	Consumed [3]string `json:"consumed" validate:"dive,tile"`
}

func NewKakan(actor int, pai string, consumed [3]string, log string) (*Kakan, error) {
	m := &Kakan{
		Action: Action{
			Message: Message{Type: TypeKakan},
			Actor:   actor,
			Log:     log,
		},
		Pai:      pai,
		Consumed: consumed,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Kakan) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeKakan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Kakan
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Kakan) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Kakan
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Kakan)(mm)
	if m.Type != TypeKakan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
