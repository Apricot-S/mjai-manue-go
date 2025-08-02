package outbound

import "github.com/Apricot-S/mjai-manue-go/internal/base"

type Dahai struct {
	action
	Pai       base.Pai
	Tsumogiri bool
}

func NewDahai(actor int, pai base.Pai, tsumogiri bool, log string) (*Dahai, error) {
	event := &Dahai{
		action: action{
			Actor: actor,
			Log:   log,
		},
		Pai:       pai,
		Tsumogiri: tsumogiri,
	}

	if err := eventValidator.Struct(event); err != nil {
		return nil, err
	}
	return event, nil
}
