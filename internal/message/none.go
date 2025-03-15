package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type None struct {
	Message
}

func NewNone() *None {
	return &None{
		Message: Message{Type: TypeNone},
	}
}

func (m *None) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type none None
	mm := (none)(*m)
	if err := json.MarshalEncode(e, &mm); err != nil {
		return err
	}
	return nil
}

func (m *None) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type none None
	var mm none
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (None)(mm)
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}
	return nil
}
