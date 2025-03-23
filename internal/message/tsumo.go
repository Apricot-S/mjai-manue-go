package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Tsumo struct {
	Action
	Pai string `json:"pai" validate:"tile"`
}

func NewTsumo(actor int, pai string, log string) (*Tsumo, error) {
	m := &Tsumo{
		Action: Action{
			Message: Message{Type: TypeTsumo},
			Actor:   actor,
			Log:     log,
		},
		Pai: pai,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Tsumo) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeTsumo {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Tsumo
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Tsumo) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Tsumo
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Tsumo)(mm)
	if m.Type != TypeTsumo {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
