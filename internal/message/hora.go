package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Hora struct {
	Action
	Target     int    `json:"target" validate:"min=0,max=3"`
	Pai        string `json:"pai" validate:"tile"`
	HoraPoints int    `json:"hora_points,omitzero" validate:"min=0"`
	Scores     []int  `json:"scores,omitempty"`
}

func NewHora(actor int, target int, pai string, horaPoints int, scores []int, log string) (*Hora, error) {
	if scores != nil && len(scores) != 4 {
		return nil, fmt.Errorf("invalid number of scores: %v", scores)
	}

	m := &Hora{
		Action: Action{
			Message: Message{Type: TypeHora},
			Actor:   actor,
			Log:     log,
		},
		Target:     target,
		Pai:        pai,
		HoraPoints: horaPoints,
		Scores:     scores,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Hora) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeHora {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner Hora
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *Hora) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner Hora
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (Hora)(mm)
	if m.Type != TypeHora {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}

	return messageValidator.Struct(m)
}
