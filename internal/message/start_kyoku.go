package message

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type StartKyoku struct {
	Message
	Bakaze     string        `json:"bakaze" validate:"wind"`
	Kyoku      int           `json:"kyoku" validate:"min=1"`
	Honba      int           `json:"honba" validate:"min=0"`
	Kyotaku    int           `json:"kyotaku" validate:"min=0"`
	Oya        int           `json:"oya" validate:"min=0,max=3"`
	DoraMarker string        `json:"dora_marker" validate:"tile"`
	Scores     []int         `json:"scores,omitempty"`
	Tehais     [4][13]string `json:"tehais" validate:"dive,dive,tile"`
}

func NewStartKyoku(
	bakaze string,
	kyoku int,
	honba int,
	kyotaku int,
	oya int,
	doraMarker string,
	scores []int,
	tehais [4][13]string,
) (*StartKyoku, error) {
	if scores != nil && len(scores) != 4 {
		return nil, fmt.Errorf("invalid number of scores: %v", scores)
	}

	m := &StartKyoku{
		Message:    Message{Type: TypeStartKyoku},
		Bakaze:     bakaze,
		Kyoku:      kyoku,
		Honba:      honba,
		Kyotaku:    kyotaku,
		Oya:        oya,
		DoraMarker: doraMarker,
		Scores:     scores,
		Tehais:     tehais,
	}

	if err := messageValidator.Struct(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *StartKyoku) MarshalJSONTo(e *jsontext.Encoder) error {
	if m.Type != TypeStartKyoku {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}
	if err := messageValidator.Struct(m); err != nil {
		return err
	}

	type inner StartKyoku
	mm := (inner)(*m)
	return json.MarshalEncode(e, &mm)
}

func (m *StartKyoku) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	type inner StartKyoku
	var mm inner
	if err := json.UnmarshalDecode(d, &mm); err != nil {
		return err
	}

	*m = (StartKyoku)(mm)
	if m.Type != TypeStartKyoku {
		return fmt.Errorf("invalid type: %v", m.Type)
	}
	if m.Scores != nil && len(m.Scores) != 4 {
		return fmt.Errorf("invalid number of scores: %v", m.Scores)
	}

	return messageValidator.Struct(m)
}
