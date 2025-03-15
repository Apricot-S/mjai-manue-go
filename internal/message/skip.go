package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// "type" is "none"
type Skip struct {
	Action
}

func NewSkip(actor int, log string) *Skip {
	return &Skip{
		Action: Action{
			Message: Message{Type: TypeNone},
			Actor:   actor,
			Log:     log,
		},
	}
}

func (m *Skip) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Skip
	mm := (inner)(*m)
	if err := json.MarshalEncode(e, &mm); err != nil {
		return err
	}
	return nil
}

func (m *Skip) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type inner Skip
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Skip)(mm)
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}
	return nil
}
