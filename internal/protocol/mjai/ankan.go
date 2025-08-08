package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Ankan struct {
	Action
	Consumed [4]string `json:"consumed" validate:"dive,tile"`
}

func NewAnkan(actor int, consumed [4]string, log string) (*Ankan, error) {
	m := &Ankan{
		Action: Action{
			Message: Message{Type: TypeAnkan},
			Actor:   actor,
			Log:     log,
		},
		Consumed: consumed,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Ankan) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeAnkan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Ankan
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Ankan) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Ankan
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Ankan)(mm)
	if m.Type != TypeAnkan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *Ankan) ToEvent() (*inbound.Ankan, error) {
	consumed := [4]base.Pai{}
	for i, c := range m.Consumed {
		p, err := base.NewPaiWithName(c)
		if err != nil {
			return nil, err
		}
		consumed[i] = *p
	}
	return inbound.NewAnkan(m.Actor, consumed)
}

func NewAnkanFromEvent(ev *outbound.Ankan) (*Ankan, error) {
	consumed := [4]string{
		ev.Consumed[0].ToString(),
		ev.Consumed[1].ToString(),
		ev.Consumed[2].ToString(),
		ev.Consumed[3].ToString(),
	}
	return NewAnkan(ev.Actor, consumed, ev.Log)
}
