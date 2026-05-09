package ai

import (
	"fmt"
	"math/rand/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
)

type ManueAgent struct {
	seed uint64
	rng  *rand.Rand
}

func NewManueAgent(seed uint64) *ManueAgent {
	agent := &ManueAgent{seed: seed}
	agent.Reset()
	return agent
}

func (a *ManueAgent) Reset() {
	a.rng = rand.New(rand.NewPCG(a.seed, 0))
}

func (a *ManueAgent) Decide(request Request) (Decision, error) {
	legalActions, err := request.Round.LegalActions(request.Self)
	if err != nil {
		return Decision{}, err
	}
	if len(legalActions) == 0 {
		return Decision{}, fmt.Errorf("cannot decide: no legal actions for player %d", request.Self.Index())
	}

	self := request.Round.Player(request.Self)
	selected, err := a.selectAction(legalActions, self)
	if err != nil {
		return Decision{}, err
	}
	return Decision{Action: selected}, nil
}

func (*ManueAgent) selectAction(legalActions []action.Action, self player.PlayerViewer) (action.Action, error) {
	if win := firstActionOfType[*action.Win](legalActions); win != nil {
		return win, nil
	}

	if self.RiichiState() == player.RiichiAccepted {
		if discard := tsumogiriDiscard(legalActions); discard != nil {
			return discard, nil
		}
		return nil, fmt.Errorf("cannot decide: no tsumogiri discard after riichi accepted")
	}

	if riichi := firstActionOfType[*action.Riichi](legalActions); riichi != nil {
		return riichi, nil
	}

	if discard := firstActionOfType[*action.Discard](legalActions); discard != nil {
		return discard, nil
	}

	if call := firstCallAction(legalActions); call != nil {
		return call, nil
	}

	if pass := firstActionOfType[*action.Pass](legalActions); pass != nil {
		return pass, nil
	}

	return legalActions[0], nil
}

func firstActionOfType[T action.Action](actions []action.Action) T {
	var zero T
	for _, a := range actions {
		if typed, ok := a.(T); ok {
			return typed
		}
	}
	return zero
}

func tsumogiriDiscard(actions []action.Action) *action.Discard {
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if ok && discard.Tsumogiri() {
			return discard
		}
	}
	return nil
}

func firstCallAction(actions []action.Action) action.Action {
	for _, a := range actions {
		switch a.(type) {
		case *action.Chii, *action.Pon, *action.CalledKan, *action.PromotedKan, *action.ConcealedKan:
			return a
		}
	}
	return nil
}
