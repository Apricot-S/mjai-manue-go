package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Join struct {
	Message
	Name string `json:"name,omitempty"`
	Room string `json:"room,omitempty"`
}

func NewJoin(name string, room string) (*Join, error) {
	m := &Join{
		Message: Message{Type: TypeJoin},
		Name:    name,
		Room:    room,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Join) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeJoin {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Join
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Join) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type inner Join
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Join)(mm)
	if m.Type != TypeJoin {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
