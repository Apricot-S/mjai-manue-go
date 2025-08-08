package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type StartGame struct {
	Message
	ID    int      `json:"id" validate:"min=0,max=3"`
	Names []string `json:"names,omitempty"`
}

func NewStartGame(id int, names []string) (*StartGame, error) {
	if names != nil && len(names) != 4 {
		return nil, fmt.Errorf("invalid number of names: %v", names)
	}

	m := &StartGame{
		Message: Message{Type: TypeStartGame},
		ID:      id,
		Names:   names,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *StartGame) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeStartGame {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Names != nil && len(m.Names) != 4 {
		return fmt.Errorf("invalid number of names: %v", m.Names)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner StartGame
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *StartGame) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner StartGame
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (StartGame)(mm)
	if m.Type != TypeStartGame {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Names != nil && len(m.Names) != 4 {
		return fmt.Errorf("invalid number of names: %v", m.Names)
	}

	return messageValidator.Struct(m)
}

func (m *StartGame) ToEvent() (*inbound.StartGame, error) {
	names := [4]string{"", "", "", ""}
	if m.Names != nil {
		names = [4]string(m.Names)
	}
	return inbound.NewStartGame(m.ID, names)
}
