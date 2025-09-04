package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

type Error struct {
	Message
}

func NewError() *Error {
	return &Error{
		Message: Message{Type: TypeError},
	}
}

func (m *Error) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeError {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Error
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Error) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Error
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Error)(mm)
	if m.Type != TypeError {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *Error) ToEvent() *inbound.Error {
	return inbound.NewError()
}
