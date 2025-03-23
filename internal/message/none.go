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

func (m *None) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner None
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *None) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner None
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (None)(mm)
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
