package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Dahai struct {
	Action
	Pai       string `json:"pai" validate:"required,tile"`
	Tsumogiri bool   `json:"tsumogiri"`
}

func NewDahai(actor int, pai string, tsumogiri bool, log string) (*Dahai, error) {
	m := &Dahai{
		Action: Action{
			Message: Message{Type: TypeDahai},
			Actor:   actor,
			Log:     log,
		},
		Pai:       pai,
		Tsumogiri: tsumogiri,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Dahai) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeDahai {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Dahai
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Dahai) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type inner Dahai
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Dahai)(mm)
	if m.Type != TypeDahai {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
