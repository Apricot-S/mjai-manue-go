package mjai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type EndKyoku struct {
	Message
}

func NewEndKyoku() *EndKyoku {
	return &EndKyoku{
		Message: Message{Type: TypeEndKyoku},
	}
}

func (m *EndKyoku) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeEndKyoku {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner EndKyoku
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *EndKyoku) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner EndKyoku
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (EndKyoku)(mm)
	if m.Type != TypeEndKyoku {
		return fmt.Errorf("invalid type: %v", m.Type)
	}

	return messageValidator.Struct(m)
}

func (m *EndKyoku) ToEvent() *inbound.EndKyoku {
	return inbound.NewEndKyoku()
}
