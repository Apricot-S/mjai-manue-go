package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
)

type TsumogiriAgent struct {
}

func NewTsumogiriAgent() *TsumogiriAgent {
	return &TsumogiriAgent{}
}

func (*TsumogiriAgent) Decide(request Request) (Decision, error) {
	legalActions, err := request.Round.LegalActions(request.Self)
	if err != nil {
		return Decision{}, err
	}
	if len(legalActions) == 0 {
		return Decision{}, fmt.Errorf("cannot decide: no legal actions for player %d", request.Self.Index())
	}

	for _, a := range legalActions {
		discard, ok := a.(*action.Discard)
		if ok && discard.Tsumogiri() {
			return Decision{Action: discard}, nil
		}

		if _, ok := a.(*action.Pass); ok {
			return Decision{Action: a}, nil
		}
	}

	// Fallback for states reached after another agent already took a non-draw action,
	// such as chii/pon/riichi declaration followed by a required discard.
	return Decision{Action: legalActions[0]}, nil
}
