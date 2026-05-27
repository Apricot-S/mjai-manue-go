package ai

import (
	"fmt"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
)

func tenpaiProb(stats TenpaiEstimatorStats, riichi bool, remainTurns int, numMelds int) float64 {
	if riichi {
		return 1.0
	}

	total, tenpai, ok := stats.YamitenCounts(remainTurns, numMelds)
	if !ok {
		return 1.0
	}
	return float64(tenpai) / float64(total)
}

// notenExhaustiveDrawTenpaiProb returns the probability that a currently
// noten player reaches tenpai before exhaustive draw, conditional on the round
// ending by exhaustive draw.
func notenExhaustiveDrawTenpaiProb(stats DrawTenpaiStats, currentTurn float64) (float64, error) {
	notenFreq := stats.ExhaustiveDrawNotenCount()

	tenpaiFreq := 0
	for turn := currentTurn + 0.25; turn <= round.FinalTurn; turn += 0.25 {
		key := strconv.FormatFloat(turn, 'f', -1, 64)
		freq, _ := stats.ExhaustiveDrawTenpaiTurnFreq(key)
		tenpaiFreq += freq
	}

	totalFreq := tenpaiFreq + notenFreq
	if totalFreq <= 0 {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw tenpai probability: frequency total must be positive")
	}
	return float64(tenpaiFreq) / float64(totalFreq), nil
}

// exhaustiveDrawTenpaiProbs adjusts current tenpai probabilities into
// exhaustive-draw tenpai probabilities. It keeps the current tenpai probability
// and adds the chance that a currently noten player reaches tenpai before
// exhaustive draw.
func exhaustiveDrawTenpaiProbs(currentTenpaiProbs [4]float64, notenTenpaiProb float64) [4]float64 {
	var probs [4]float64
	for playerID, currentTenpaiProb := range currentTenpaiProbs {
		probs[playerID] = currentTenpaiProb + (1.0-currentTenpaiProb)*notenTenpaiProb
	}
	return probs
}

// ryukyokuScoreDelta returns the score change vector for exhaustive draw
// tenpai payments.
func ryukyokuScoreDelta(tenpais [4]bool) scoreDelta {
	points := service.RyukyokuPoints(tenpais)
	var delta scoreDelta
	for i, point := range points {
		delta[i] = float64(point)
	}
	return delta
}

// exhaustiveDrawScoreDeltaDistFromTenpaiProbs returns the score change
// distribution assuming the round ends in an exhaustive draw.
func exhaustiveDrawScoreDeltaDistFromTenpaiProbs(tenpaiProbs [4]float64) scoreDeltaProbDist {
	tenpaisDist := aheadVectorProbDist{{}: 1.0}
	for playerID, tenpaiProb := range tenpaiProbs {
		var tenpais aheadVector
		tenpais[playerID] = 1
		tenpaisDist = addAheadVectorProbDists(tenpaisDist, newAheadVectorProbDist(map[aheadVector]float64{
			{}:      1.0 - tenpaiProb,
			tenpais: tenpaiProb,
		}))
	}

	return tenpaisDist.mapValueScoreDelta(func(tenpais aheadVector) scoreDelta {
		return ryukyokuScoreDelta(aheadVectorToBoolArray(tenpais))
	})
}

// futureExhaustiveDrawScoreDeltaDist returns the exhaustive-draw score change
// distribution from current tenpai probabilities. It first adjusts current
// tenpai probabilities by the chance that currently noten players reach tenpai
// before exhaustive draw.
func futureExhaustiveDrawScoreDeltaDist(
	currentTenpaiProbs [4]float64,
	notenTenpaiProb float64,
) scoreDeltaProbDist {
	return exhaustiveDrawScoreDeltaDistFromTenpaiProbs(exhaustiveDrawTenpaiProbs(currentTenpaiProbs, notenTenpaiProb))
}

// exhaustiveDrawAvgPts returns self's expected score change assuming the round
// ends in an exhaustive draw and tenpaiProbs already represent exhaustive-draw
// tenpai probabilities.
func exhaustiveDrawAvgPts(selfID int, tenpaiProbs [4]float64) float64 {
	return exhaustiveDrawScoreDeltaDistFromTenpaiProbs(tenpaiProbs).expected()[selfID]
}

func aheadVectorToBoolArray(value aheadVector) [4]bool {
	var result [4]bool
	for i, v := range value {
		result[i] = v != 0
	}
	return result
}
