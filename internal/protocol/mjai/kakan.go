package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type Kakan struct {
	Action
	Pai      string    `json:"pai" validate:"tile"`
	Consumed [3]string `json:"consumed" validate:"dive,tile"`
}

func NewKakan(actor int, pai string, consumed [3]string, log string) (*Kakan, error) {
	m := &Kakan{
		Action: Action{
			Message: Message{Type: TypeKakan},
			Actor:   actor,
			Log:     log,
		},
		Pai:      pai,
		Consumed: consumed,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Kakan) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeKakan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Kakan
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Kakan) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Kakan
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Kakan)(mm)
	if m.Type != TypeKakan {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *Kakan) ToEvent() (*inbound.Kakan, error) {
	added, err := base.NewPaiWithName(m.Pai)
	if err != nil {
		return nil, err
	}

	consumed := [2]base.Pai{}
	for i, c := range m.Consumed[:2] {
		p, err := base.NewPaiWithName(c)
		if err != nil {
			return nil, err
		}
		consumed[i] = *p
	}

	// Heuristic: The last tile in the consumed is considered `taken`.
	// This is a simplification and may not reflect the actual game state.
	taken, err := base.NewPaiWithName(m.Consumed[2])
	if err != nil {
		return nil, err
	}

	// Target is temporarily set to a value that does not overlap with Actor.
	// There is no problem when updating the Player because the Target is obtained from the Pon.
	target := (3 + m.Actor) % 4

	// Target is not provided in event data
	return inbound.NewKakan(m.Actor, target, *taken, consumed, *added)
}

func NewKakanFromEvent(ev *outbound.Kakan) (*Kakan, error) {
	consumed := [3]string{
		ev.Consumed[0].ToString(),
		ev.Consumed[1].ToString(),
		ev.Taken.ToString(),
	}
	return NewKakan(ev.Actor, ev.Added.ToString(), consumed, ev.Log)
}
