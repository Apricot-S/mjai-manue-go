package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Hello struct {
	Message
	Protocol        string `json:"protocol,omitempty"`
	ProtocolVersion int    `json:"protocol_version,omitzero" validate:"min=0"`
}

func NewHello(protocol string, protocolVersion int) (*Hello, error) {
	m := &Hello{
		Message:         Message{Type: TypeHello},
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Hello) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeHello {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Hello
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Hello) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type inner Hello
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Hello)(mm)
	if m.Type != TypeHello {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
