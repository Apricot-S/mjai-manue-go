package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/ai/core"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

func (a *ManueAI) getScoreChangesDistOnHora(
	state game.StateViewer,
	playerID int,
	horaPointsDist *core.ProbDist[[]float64],
) *core.ProbDist[[]float64] {
	tsumoHoraProb := float64(a.stats.NumTsumoHoras) / float64(a.stats.NumHoras)
	unitDistMap := core.NewHashMap[[]float64]()

	for _, target := range state.Players() {
		var changes []float64
		if target.ID() != playerID {
			changes = []float64{0.0, 0.0, 0.0, 0.0}
			changes[playerID] = 1
			changes[target.ID()] = -1
		} else if playerID == state.Oya().ID() {
			changes = []float64{-1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0}
			changes[playerID] = 1
		} else {
			changes = []float64{-1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0}
			changes[playerID] = 1
			changes[state.Oya().ID()] = -1.0 / 2.0
		}

		var prob float64
		if target.ID() == playerID {
			prob = tsumoHoraProb
		} else {
			prob = (1.0 - tsumoHoraProb) / 3.0
		}
		unitDistMap.Set(changes, prob)
	}

	u := core.NewProbDist(unitDistMap)
	return core.Mult[[]float64, []float64, []float64](horaPointsDist, u)
}
