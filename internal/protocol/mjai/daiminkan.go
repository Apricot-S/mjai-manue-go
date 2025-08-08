package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Daiminkan struct {
	Action
	Target   int       `json:"target" validate:"min=0,max=3"`
	Pai      string    `json:"pai" validate:"tile"`
	Consumed [3]string `json:"consumed" validate:"dive,tile"`
}

func NewDaiminkan(actor int, target int, pai string, consumed [3]string, log string) (*Daiminkan, error) {
	m := &Daiminkan{
		Action: Action{
			Message: Message{Type: TypeDaiminkan},
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

func (m *Daiminkan) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeDaiminkan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Daiminkan
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Daiminkan) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Daiminkan
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Daiminkan)(mm)
	if m.Type != TypeDaiminkan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *Daiminkan) ToEvent() (*inbound.Daiminkan, error) {
	taken, err := base.NewPaiWithName(m.Pai)
	if err != nil {
		return nil, err
	}

	consumed := [3]base.Pai{}
	for i, c := range m.Consumed {
		p, err := base.NewPaiWithName(c)
		if err != nil {
			return nil, err
		}
		consumed[i] = *p
	}

	return inbound.NewDaiminkan(m.Actor, m.Target, *taken, consumed)
}

func NewDaiminkanFromEvent(ev *outbound.Daiminkan) (*Daiminkan, error) {
	consumed := [3]string{ev.Consumed[0].ToString(), ev.Consumed[1].ToString(), ev.Consumed[2].ToString()}
	return NewDaiminkan(ev.Actor, ev.Target, ev.Taken.ToString(), consumed, ev.Log)
}
