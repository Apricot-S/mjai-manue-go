package ai

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"

type TsumogiriAgent struct {
}

func NewTsumogiriAgent() *TsumogiriAgent {
	return &TsumogiriAgent{}
}

func (*TsumogiriAgent) Decide(request Request) (Decision, error) {
	p := request.Round.Player(request.Self)
	drawnTile := p.DrawnTile()
	if drawnTile == nil {
		return Decision{Action: action.NewPass(request.Self)}, nil
	}

	a, err := action.NewDiscard(request.Self, *drawnTile, true)
	if err != nil {
		return Decision{}, err
	}
	return Decision{Action: a}, nil
}
