package ai

import (
	"fmt"
	"math/rand/v2"
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

func (*ManueAgent) Decide(request Request) (Decision, error) {
	legalActions, err := request.Round.LegalActions(request.Self)
	if err != nil {
		return Decision{}, err
	}
	if len(legalActions) == 0 {
		return Decision{}, fmt.Errorf("cannot decide: no legal actions for player %d", request.Self.Index())
	}

	return Decision{Action: legalActions[0]}, nil
}
