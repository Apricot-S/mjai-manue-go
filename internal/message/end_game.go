package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type EndGame struct {
	Message
}

func NewEndGame() *EndGame {
	return &EndGame{
		Message: Message{Type: TypeEndGame},
	}
}

func (m *EndGame) MarshalJSONTo(e *jsontext.Encoder, opts jsontext.Options) error {
	if m.Type != TypeEndGame {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner EndGame
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *EndGame) UnmarshalJSONFrom(d *jsontext.Decoder, opts jsontext.Options) error {
	type inner EndGame
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (EndGame)(mm)
	if m.Type != TypeEndGame {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}
