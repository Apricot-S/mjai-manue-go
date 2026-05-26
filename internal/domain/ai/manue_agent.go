package ai

import (
	"fmt"
	"math/rand/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type ManueAgent struct {
	seed uint64
	rng  *rand.Rand
	deps ManueAgentDeps
}

func NewManueAgent(seed uint64, deps ManueAgentDeps) (*ManueAgent, error) {
	if deps.Stats == nil {
		return nil, fmt.Errorf("cannot create ManueAgent: stats dependency is required")
	}
	if err := validateManueStats(deps.Stats); err != nil {
		return nil, fmt.Errorf("cannot create ManueAgent: %w", err)
	}
	if deps.Danger == nil {
		return nil, fmt.Errorf("cannot create ManueAgent: danger estimator dependency is required")
	}
	agent := &ManueAgent{
		seed: seed,
		deps: deps,
	}
	agent.Reset()
	return agent, nil
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

	if win := firstActionOfType[*action.Win](legalActions); win != nil {
		// Always take a winning action when it is legal.
		// The current policy does not allow passing on win opportunities.
		return Decision{Action: win}, nil
	}

	if request.Round.Player(request.Self).CanDiscard() {
		return a.decideSelfTurn(legalActions, request.Round, request.Self)
	}
	return a.decideOtherDiscardReaction(legalActions, request.Round, request.Self)
}

func (a *ManueAgent) decideSelfTurn(
	legalActions []action.Action,
	state round.StateViewer,
	selfSeat seat.Seat,
) (Decision, error) {
	if state == nil {
		return Decision{}, fmt.Errorf("cannot decide self turn: state is required")
	}
	self := state.Player(selfSeat)
	if self == nil {
		return Decision{}, fmt.Errorf("cannot decide self turn: self player is required")
	}

	if self.RiichiState() == player.RiichiAccepted {
		// After riichi is accepted, always discard the drawn tile.
		// Concealed kan is intentionally ignored even if it is legal.
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
	candidates, err = a.evaluateActionCandidates(state, selfSeat, candidates)
	if err != nil {
		return Decision{}, err
	}

	candidate := chooseBestCandidate(candidates, true)
	log := formatCandidateLog(candidates)
	return Decision{
		Action: candidate.action,
		Log:    log,
		Trace:  formatDecisionTrace(log, &candidate),
	}, nil
}

func (a *ManueAgent) decideOtherDiscardReaction(
	legalActions []action.Action,
	state round.StateViewer,
	selfSeat seat.Seat,
) (Decision, error) {
	if state == nil {
		return Decision{}, fmt.Errorf("cannot decide other discard reaction: state is required")
	}
	self := state.Player(selfSeat)
	if self == nil {
		return Decision{}, fmt.Errorf("cannot decide other discard reaction: self player is required")
	}

	candidates, err := getOtherDiscardReactionCandidates(legalActions, self)
	if err != nil {
		return Decision{}, err
	}
	if len(candidates) == 0 {
		return Decision{}, fmt.Errorf("cannot decide other discard reaction: no reaction candidate")
	}
	candidates, err = a.evaluateActionCandidates(state, selfSeat, candidates)
	if err != nil {
		return Decision{}, err
	}

	candidate := chooseBestCandidate(candidates, false)
	log := formatCandidateLog(candidates)
	return Decision{
		Action: candidate.action,
		Log:    log,
		Trace:  formatDecisionTrace(log, &candidate),
	}, nil
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
