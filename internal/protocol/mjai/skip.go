package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// "type" is "none"
type Skip struct {
	Action
}

func NewSkip(actor int, log string) (*Skip, error) {
	m := &Skip{
		Action: Action{
			Message: Message{Type: TypeNone},
			Actor:   actor,
			Log:     log,
		},
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Skip) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Skip
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Skip) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Skip
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Skip)(mm)
	if m.Type != TypeNone {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func NewSkipFromEvent(ev *outbound.Skip) (*Skip, error) {
	return NewSkip(ev.Actor, ev.Log)
}
