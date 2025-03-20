package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Dora struct {
	Message
	DoraMarker string `json:"dora_marker" validate:"required,tile"`
}

func NewDora(doraMarker string) (*Dora, error) {
	m := &Dora{
		Message:    Message{Type: TypeDora},
		DoraMarker: doraMarker,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Dora) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeDora {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Dora
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Dora) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type inner Dora
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Dora)(mm)
	if m.Type != TypeDora {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
