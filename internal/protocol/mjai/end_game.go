package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

type EndGame struct {
	Message
}

func NewEndGame() *EndGame {
	return &EndGame{
		Message: Message{Type: TypeEndGame},
	}
}

func (m *EndGame) MarshalJSONTo(e *jsontext.Encoder) error {
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

func (m *EndGame) UnmarshalJSONFrom(d *jsontext.Decoder) error {
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

func (m *EndGame) ToEvent() *inbound.EndGame {
	return inbound.NewEndGame()
}
