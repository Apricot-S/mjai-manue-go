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
	deps ManueAgentDeps
}

type ManueAgentDeps struct {
	Stats ManueStats
}

type ManueStats interface {
	WinScoreStats
}

type WinScoreStats interface {
	NumWins() int
	NumSelfDrawWins() int
	NonDealerWinPointFreqs() map[string]int
	DealerWinPointFreqs() map[string]int
}

func NewManueAgent(seed uint64) *ManueAgent {
	return newManueAgent(seed, ManueAgentDeps{})
}

func NewManueAgentWithDeps(seed uint64, deps ManueAgentDeps) *ManueAgent {
	return newManueAgent(seed, deps)
}

func newManueAgent(seed uint64, deps ManueAgentDeps) *ManueAgent {
	agent := &ManueAgent{
		seed: seed,
		deps: deps,
	}
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
	if win := firstActionOfType[*action.Win](legalActions); win != nil {
		return Decision{Action: win}, nil
	}
	if self.CanDiscard() {
		return a.decideSelfTurn(legalActions, self)
	}
	return a.decideOtherDiscardReaction(legalActions)
}

func (*ManueAgent) decideSelfTurn(legalActions []action.Action, self player.PlayerViewer) (Decision, error) {
	if self.RiichiState() == player.RiichiAccepted {
		if discard := tsumogiriDiscard(legalActions); discard != nil {
			return Decision{Action: discard}, nil
		}
		return Decision{}, fmt.Errorf("cannot decide: no tsumogiri discard after riichi accepted")
	}

	candidates, err := getSelfTurnCandidates(legalActions, self)
	if err != nil {
		return Decision{}, err
	}
	if len(candidates) == 0 {
		return Decision{}, fmt.Errorf("cannot decide self turn: no self-turn candidate")
	}

	candidate := chooseBestCandidate(candidates, true)
	log := formatCandidateLog(candidates)
	return Decision{
		Action: candidate.action,
		Log:    log,
		Trace:  formatDecisionTrace(log, &candidate),
	}, nil
}

func (*ManueAgent) decideOtherDiscardReaction(legalActions []action.Action) (Decision, error) {
	if call := firstCallAction(legalActions); call != nil {
		return Decision{Action: call}, nil
	}

	if pass := firstActionOfType[*action.Pass](legalActions); pass != nil {
		return Decision{Action: pass}, nil
	}

	return Decision{}, fmt.Errorf("cannot decide other discard reaction: no call or pass candidate")
}

func firstActionOfType[T action.Action](actions []action.Action) T {
	for _, a := range actions {
		if typed, ok := a.(T); ok {
			return typed
		}
	}
	var zero T
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
