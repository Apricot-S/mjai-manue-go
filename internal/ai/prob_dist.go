package ai

import (
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/ai/core"
	"github.com/Apricot-S/mjai-manue-go/internal/ai/estimator"
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

func (a *ManueAI) getRyukyokuAveragePoints(
	state game.StateViewer,
	playerID int,
	selfTenpai bool,
) float64 {
	notenRyukyokuTenpaiProb := a.getNotenRyukyokuTenpaiProb(state)
	var ryukyokuTenpaiProbs [4]float64
	for i := range 4 {
		player := &state.Players()[i]
		var currentTenpaiProb = 0.0
		if player.ID() == playerID {
			if selfTenpai {
				currentTenpaiProb = 1.0
			} else {
				currentTenpaiProb = 0.0
			}
		} else {
			currentTenpaiProb = a.tenpaiProbEstimator.Estimate(player, state)
		}
		ryukyokuTenpaiProbs[i] = currentTenpaiProb + (1.0-currentTenpaiProb)*notenRyukyokuTenpaiProb
	}

	result := 0.0
	for i := range 1 << 4 {
		var tenpais [4]bool
		for j := range 4 {
			tenpais[j] = (i & (1 << j)) != 0
		}
		prob := 1.0
		numTenpais := 0
		for j := range 4 {
			if tenpais[j] {
				prob *= ryukyokuTenpaiProbs[j]
				numTenpais++
			} else {
				prob *= 1.0 - ryukyokuTenpaiProbs[j]
			}
		}
		if prob > 0.0 {
			var points float64
			if tenpais[playerID] {
				if numTenpais == 4 {
					points = 0.0
				} else {
					points = 3000.0 / float64(numTenpais)
				}
			} else {
				if numTenpais == 0 {
					points = 0.0
				} else {
					points = -3000.0 / float64(4-numTenpais)
				}
			}
			result += prob * points
		}
	}
	return result
}

// Distribution of score changes assuming the kyoku ends with ryukyoku.
func (a *ManueAI) getScoreChangesDistOnRyukyoku(
	state game.StateViewer,
	playerID int,
	selfTenpai bool,
) *core.ProbDist[[]float64] {
	notenRyukyokuTenpaiProb := a.getNotenRyukyokuTenpaiProb(state)
	hm1 := core.NewHashMap[[]float64]()
	hm1.Set([]float64{0.0, 0.0, 0.0, 0.0}, 1.0)
	tenpaisDist := core.NewProbDist(hm1)

	for _, player := range state.Players() {
		var currentTenpaiProb = 0.0
		if player.ID() == playerID {
			if selfTenpai {
				currentTenpaiProb = 1.0
			} else {
				currentTenpaiProb = 0.0
			}
		} else {
			currentTenpaiProb = a.tenpaiProbEstimator.Estimate(&player, state)
		}
		ryukyokuTenpaiProb := currentTenpaiProb + (1.0-currentTenpaiProb)*notenRyukyokuTenpaiProb

		tenpais := make([]float64, 4)
		for i := range 4 {
			if player.ID() == playerID {
				tenpais[i] = 1.0
			} else {
				tenpais[i] = 0.0
			}
		}

		hm2 := core.NewHashMap[[]float64]()
		hm2.Set([]float64{0.0, 0.0, 0.0, 0.0}, 1.0-ryukyokuTenpaiProb)
		hm2.Set(tenpais, ryukyokuTenpaiProb)
		dist := core.NewProbDist(hm2)
		tenpaisDist = core.Add[[]float64, []float64, []float64](tenpaisDist, dist)
	}

	return tenpaisDist.MapValue(tenpaisToRyukyokuPointsFloat)
}

func tenpaisToRyukyokuPointsFloat(tenpais []float64) []float64 {
	t := [4]bool{}
	for i := range tenpais {
		t[i] = tenpais[i] != 0.0
	}
	r := game.TenpaisToRyukyokuPoints(t)
	ret := make([]float64, 4)
	for i := range r {
		ret[i] = float64(r[i])
	}
	return ret
}

// Probability that the player is tenpai at the end of the kyoku if
// the player is currently noten and the kyoku ends with ryukyoku.
func (a *ManueAI) getNotenRyukyokuTenpaiProb(state game.StateViewer) float64 {
	notenFreq := float64(a.stats.RyukyokuTenpaiStat.Noten)
	tenpaiFreq := 0.0
	t := float64(state.Turn()) + 1.0/4.0
	for t <= float64(game.FinalTurn) {
		n := strconv.FormatFloat(t, 'f', -1, 64)
		tenpaiFreq += float64(a.stats.RyukyokuTenpaiStat.TenpaiTurnDistribution[n])
		t += 1.0 / 4.0
	}
	return tenpaiFreq / (tenpaiFreq + notenFreq)
}

func (a *ManueAI) getSafeProbs(
	state game.StateViewer,
	playerID int,
	dahaiCandidates []game.Pai,
) (map[string]float64, error) {
	safeProbs := make(map[string]float64, len(dahaiCandidates))
	for _, pai := range dahaiCandidates {
		var key string
		if pai.IsUnknown() {
			key = pai.ToString()
		} else {
			key = "none"
		}
		safeProbs[key] = 1.0
	}

	me := &state.Players()[playerID]
	for _, player := range state.Players() {
		if player.ID() == playerID {
			continue
		}

		scene, err := estimator.NewScene(state, me, &player)
		if err != nil {
			return nil, err
		}
		tenpaiProb := a.tenpaiProbEstimator.Estimate(&player, state)
		for _, pai := range dahaiCandidates {
			if pai.IsUnknown() {
				continue
			}

			isAnpai, err := scene.Evaluate("anpai", &pai)
			if err != nil {
				return nil, err
			}
			if isAnpai {
				continue
			}

			probInfo, err := a.dangerEstimator.EstimateProb(scene, &pai)
			if err != nil {
				return nil, err
			}
			safeProb := 1.0 - tenpaiProb*probInfo.Prob
			safeProbs[pai.ToString()] *= safeProb
		}
	}

	return safeProbs, nil
}

func (a *ManueAI) getRyukyokuProb(state game.StateViewer) float64 {
	currentTurn := state.Turn()
	den := 0.0
	for _, t := range a.stats.NumTurnsDistribution[currentTurn:] {
		den += t
	}
	return a.stats.RyukyokuRatio / den
}
