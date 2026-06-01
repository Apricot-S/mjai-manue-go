package ai

import (
	"fmt"
	"math"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
)

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
	return math.Pow(prob, noSelfWinExhaustiveDrawExponent()), nil
}

func noSelfWinExhaustiveDrawExponent() float64 {
	return float64(common.NumPlayers-1) / float64(common.NumPlayers)
}

func expectedRemainingTurns(stats RoundEndStats, currentTurn float64) (int, error) {
	if currentTurn < 0 || currentTurn > round.FinalTurn {
		return 0, fmt.Errorf("cannot estimate expected remaining turns: current turn is out of range")
	}

	currentTurnIndex := int(math.Round(currentTurn))
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

// remainingRoundEndProbs splits the probability mass left after accounting for
// self's future win into exhaustive draw and other players' win.
func remainingRoundEndProbs(selfWinProb float64, exhaustiveDrawProbOnSelfNoWin float64) (exhaustiveDrawProb float64, otherWinProb float64, err error) {
	if selfWinProb < 0.0 || selfWinProb > 1.0 {
		return 0.0, 0.0, fmt.Errorf("cannot estimate remaining round-end probabilities: self win probability must be between 0 and 1")
	}
	if exhaustiveDrawProbOnSelfNoWin < 0.0 || exhaustiveDrawProbOnSelfNoWin > 1.0 {
		return 0.0, 0.0, fmt.Errorf("cannot estimate remaining round-end probabilities: exhaustive-draw probability must be between 0 and 1")
	}

	noSelfWinProb := 1.0 - selfWinProb
	return noSelfWinProb * exhaustiveDrawProbOnSelfNoWin,
		noSelfWinProb * (1.0 - exhaustiveDrawProbOnSelfNoWin),
		nil
}
