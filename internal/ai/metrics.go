package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type metric struct {
}

type metrics map[string]metric

func (a *ManueAI) getMetrics(
	state game.StateViewer,
	playerID int,
	dahaiCandidates []game.Pai,
	reachDahaiCandidates []game.Pai,
	forbiddenDahais []game.Pai,
) (pai *game.Pai, isReach bool, err error) {
	// TODO: Implement logic.
	return &dahaiCandidates[len(dahaiCandidates)-1], false, nil
}
