package ai

import (
	"fmt"
	"math"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
)

// winScoreFactor returns how one win point unit changes all players' scores.
//
// actorID is the winner. targetID is the winner for self-draw wins, or the
// discarder for ron wins. dealerID is the round dealer.
func winScoreFactor(actorID int, targetID int, dealerID int) scoreDelta {
	if targetID != actorID {
		// Ron: the discarder pays the full win points.
		var factor scoreDelta
		factor[actorID] = 1.0
		factor[targetID] = -1.0
		return factor
	}

	if actorID == dealerID {
		// Dealer self-draw: each non-dealer pays one third.
		factor := scoreDelta{-1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0}
		factor[actorID] = 1.0
		return factor
	}

	// Non-dealer self-draw: the dealer pays half, each other non-dealer pays a quarter.
	factor := scoreDelta{-1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0}
	factor[actorID] = 1.0
	factor[dealerID] = -1.0 / 2.0
	return factor
}

// winScoreFactorDist returns the distribution of score factors for a winner.
//
// selfDrawProb is the probability that the win is by self draw. Ron targets are
// treated as uniformly distributed among the three other players.
func winScoreFactorDist(actorID int, dealerID int, selfDrawProb float64) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, 4)
	ronTargetProb := (1.0 - selfDrawProb) / 3.0
	for targetID := range 4 {
		var prob float64
		if targetID == actorID {
			prob = selfDrawProb
		} else {
			prob = ronTargetProb
		}
		dist[winScoreFactor(actorID, targetID, dealerID)] = prob
	}
	return newScoreDeltaProbDist(dist)
}

// winPointsDist returns a probability distribution from win-points frequencies.
func winPointsDist(pointFreqs map[string]int) (scalarProbDist, error) {
	totalFreqs := pointFreqs["total"]
	if totalFreqs <= 0 {
		return nil, fmt.Errorf("invalid win points frequencies: total must be positive")
	}
	totalFreqsFloat := float64(totalFreqs)

	dist := make(map[float64]float64, len(pointFreqs)-1)
	for points, freq := range pointFreqs {
		if points == "total" {
			continue
		}
		parsedPoints, err := strconv.ParseFloat(points, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid win points frequency key %q: %w", points, err)
		}
		dist[parsedPoints] = float64(freq) / totalFreqsFloat
	}
	return newScalarProbDist(dist), nil
}

// randomWinScoreDeltaDist returns the score-change distribution for a random
// win by actorID.
func randomWinScoreDeltaDist(
	actorID int,
	dealerID int,
	selfDrawProb float64,
	pointFreqs map[string]int,
) (scoreDeltaProbDist, error) {
	pointsDist, err := winPointsDist(pointFreqs)
	if err != nil {
		return nil, err
	}
	return multiplyScalarScoreDeltaProbDists(
		pointsDist,
		winScoreFactorDist(actorID, dealerID, selfDrawProb),
	), nil
}

func randomWinScoreDeltaDistFromStats(
	actorID int,
	dealerID int,
	stats WinScoreStats,
) (scoreDeltaProbDist, error) {
	if stats.NumWins() <= 0 {
		return nil, fmt.Errorf("cannot build random win score delta distribution: numWins must be positive")
	}

	pointFreqs := stats.NonDealerWinPointFreqs()
	if actorID == dealerID {
		pointFreqs = stats.DealerWinPointFreqs()
	}
	return randomWinScoreDeltaDist(
		actorID,
		dealerID,
		float64(stats.NumSelfDrawWins())/float64(stats.NumWins()),
		pointFreqs,
	)
}

func exhaustiveDrawProb(stats RoundEndStats, currentTurn float64) (float64, error) {
	currentTurnIndex := int(currentTurn)
	turnDistribution := stats.TurnDistribution()
	if currentTurnIndex < 0 || currentTurnIndex >= numTurnDistributionEntries {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw probability: current turn is out of range")
	}

	remainingProb := 0.0
	for _, prob := range turnDistribution[currentTurnIndex:] {
		remainingProb += prob
	}
	if remainingProb <= 0 {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw probability: remaining turn probability must be positive")
	}
	return stats.ExhaustiveDrawRatio() / remainingProb, nil
}

func exhaustiveDrawProbOnSelfNoWin(stats RoundEndStats, currentTurn float64) (float64, error) {
	prob, err := exhaustiveDrawProb(stats, currentTurn)
	if err != nil {
		return 0, err
	}
	return math.Pow(prob, 3.0/4.0), nil
}

func expectedRemainingTurns(stats RoundEndStats, currentTurn float64) (int, error) {
	currentTurnIndex := int(math.Round(currentTurn))
	if currentTurnIndex < 0 || currentTurnIndex >= numTurnDistributionEntries {
		return 0, fmt.Errorf("cannot estimate expected remaining turns: current turn is out of range")
	}

	turnDistribution := stats.TurnDistribution()
	num := 0.0
	den := 0.0
	for i := currentTurnIndex; i < numTurnDistributionEntries; i++ {
		prob := turnDistribution[i]
		num += prob * (float64(i) - math.Round(currentTurn) + 0.5)
		den += prob
	}
	if den == 0 {
		return 0, nil
	}
	return int(math.Round(num / den)), nil
}

func tenpaiProb(stats TenpaiEstimatorStats, riichi bool, remainTurns int, numMelds int) (float64, error) {
	if riichi {
		return 1.0, nil
	}

	total, tenpai, ok := stats.YamitenCounts(remainTurns, numMelds)
	if !ok {
		return 1.0, nil
	}
	if total <= 0 {
		return 0, fmt.Errorf("cannot estimate tenpai probability: yamiten total must be positive")
	}
	if tenpai < 0 || tenpai > total {
		return 0, fmt.Errorf("cannot estimate tenpai probability: yamiten tenpai must be between 0 and total")
	}
	return float64(tenpai) / float64(total), nil
}

// notenExhaustiveDrawTenpaiProb returns the probability that a currently
// noten player reaches tenpai before exhaustive draw, conditional on the round
// ending by exhaustive draw.
func notenExhaustiveDrawTenpaiProb(stats DrawTenpaiStats, currentTurn float64) (float64, error) {
	notenFreq := stats.ExhaustiveDrawNotenCount()
	if notenFreq < 0 {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw tenpai probability: noten count must be non-negative")
	}

	tenpaiFreq := 0
	for turn := currentTurn + 0.25; turn <= round.FinalTurn; turn += 0.25 {
		key := strconv.FormatFloat(turn, 'f', -1, 64)
		freq, ok := stats.ExhaustiveDrawTenpaiTurnFreq(key)
		if !ok {
			return 0, fmt.Errorf("cannot estimate exhaustive-draw tenpai probability: missing tenpai turn frequency for turn %s", key)
		}
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

// ryukyokuScoreDeltaDist returns the score change distribution assuming the
// round ends in an exhaustive draw.
func ryukyokuScoreDeltaDist(tenpaiProbs [4]float64) scoreDeltaProbDist {
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

func exhaustiveDrawScoreDeltaDist(
	currentTenpaiProbs [4]float64,
	notenTenpaiProb float64,
) scoreDeltaProbDist {
	return ryukyokuScoreDeltaDist(exhaustiveDrawTenpaiProbs(currentTenpaiProbs, notenTenpaiProb))
}

func exhaustiveDrawAvgPts(selfID int, tenpaiProbs [4]float64) float64 {
	return ryukyokuScoreDeltaDist(tenpaiProbs).expected()[selfID]
}

func aheadVectorToBoolArray(value aheadVector) [4]bool {
	var result [4]bool
	for i, v := range value {
		result[i] = v != 0
	}
	return result
}
