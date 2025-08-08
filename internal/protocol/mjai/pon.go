package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Pon struct {
	Action
	Target   int       `json:"target" validate:"min=0,max=3"`
	Pai      string    `json:"pai" validate:"tile"`
	Consumed [2]string `json:"consumed" validate:"dive,tile"`
}

func NewPon(actor int, target int, pai string, consumed [2]string, log string) (*Pon, error) {
	m := &Pon{
		Action: Action{
			Message: Message{Type: TypePon},
			Actor:   actor,
			Log:     log,
		},
		Target:   target,
		Pai:      pai,
		Consumed: consumed,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Pon) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypePon {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Pon
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Pon) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Pon
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Pon)(mm)
	if m.Type != TypePon {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *Pon) ToEvent() (*inbound.Pon, error) {
	taken, err := base.NewPaiWithName(m.Pai)
	if err != nil {
		return nil, err
	}

	consumed := [2]base.Pai{}
	for i, c := range m.Consumed {
		p, err := base.NewPaiWithName(c)
		if err != nil {
			return nil, err
		}
		consumed[i] = *p
	}

	return inbound.NewPon(m.Actor, m.Target, *taken, consumed)
}

func NewPonFromEvent(ev *outbound.Pon) (*Pon, error) {
	consumed := [2]string{ev.Consumed[0].ToString(), ev.Consumed[1].ToString()}
	return NewPon(ev.Actor, ev.Target, ev.Taken.ToString(), consumed, ev.Log)
}
