package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/base"

type Tsumo struct {
	action
	Pai base.Pai
}

func NewTsumo(actor int, pai base.Pai, log string) (*Tsumo, error) {
	event := &Tsumo{
		action: action{
			Actor: actor,
			Log:   log,
		},
		Pai: pai,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
