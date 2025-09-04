package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type Dahai struct {
	Action
	Pai       string `json:"pai" validate:"tile"`
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

func (m *Dahai) MarshalJSONTo(e *jsontext.Encoder) error {
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

func (m *Dahai) UnmarshalJSONFrom(d *jsontext.Decoder) error {
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

func (m *Dahai) ToEvent() (*inbound.Dahai, error) {
	pai, err := base.NewPaiWithName(m.Pai)
	if err != nil {
		return nil, err
	}

	return inbound.NewDahai(m.Actor, *pai, m.Tsumogiri)
}

func NewDahaiFromEvent(ev *outbound.Dahai) (*Dahai, error) {
	return NewDahai(ev.Actor, ev.Pai.ToString(), ev.Tsumogiri, ev.Log)
}
