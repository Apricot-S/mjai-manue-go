package mjai

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
)

type Dora struct {
	Message
	DoraMarker string `json:"dora_marker" validate:"tile"`
}

func NewDora(doraMarker string) (*Dora, error) {
	m := &Dora{
		Message:    Message{Type: TypeDora},
		DoraMarker: doraMarker,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Dora) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeDora {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Dora
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Dora) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Dora
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Dora)(mm)
	if m.Type != TypeDora {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *Dora) ToEvent() (*inbound.Dora, error) {
	doraMarker, err := base.NewPaiWithName(m.DoraMarker)
	if err != nil {
		return nil, err
	}

	return inbound.NewDora(*doraMarker)
}
