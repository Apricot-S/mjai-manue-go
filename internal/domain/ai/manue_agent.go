package ai

import (
	"fmt"
	"math/rand/v2"
)

type ManueAgent struct {
	rng *rand.Rand
}

func NewManueAgent(seed uint64) *ManueAgent {
	return &ManueAgent{
		rng: rand.New(rand.NewPCG(seed, 0)),
	}
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
