package ai

import (
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
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
		var factor scoreDelta
		for id := range factor {
			factor[id] = -dealerSelfDrawPaymentFactor()
		}
		factor[actorID] = 1.0
		return factor
	}

	// Non-dealer self-draw: the dealer pays half, each other non-dealer pays a quarter.
	var factor scoreDelta
	for id := range factor {
		factor[id] = -nonDealerSelfDrawOtherPaymentFactor()
	}
	factor[actorID] = 1.0
	factor[dealerID] = -1.0 / 2.0
	return factor
}

func dealerSelfDrawPaymentFactor() float64 {
	return 1.0 / float64(common.NumPlayers-1)
}

func nonDealerSelfDrawOtherPaymentFactor() float64 {
	return 1.0 / (2.0 * float64(common.NumPlayers-2))
}

// winScoreFactorDist returns the distribution of score factors for a winner.
//
// selfDrawProb is the probability that the win is by self draw. Ron targets are
// treated as uniformly distributed among the other players.
func winScoreFactorDist(actorID int, dealerID int, selfDrawProb float64) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, common.NumPlayers)
	ronTargetProb := (1.0 - selfDrawProb) / float64(common.NumPlayers-1)
	for targetID := range common.NumPlayers {
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

func winPointsDistFromValidatedStats(pointFreqs map[string]int) scalarProbDist {
	totalFreqsFloat := float64(pointFreqs["total"])
	dist := make(map[float64]float64, len(pointFreqs)-1)
	for points, freq := range pointFreqs {
		if points == "total" {
			continue
		}
		parsedPoints, _ := strconv.ParseFloat(points, 64)
		dist[parsedPoints] = float64(freq) / totalFreqsFloat
	}
	return newScalarProbDist(dist)
}

func randomWinScoreDeltaDistFromStats(
	actorID int,
	dealerID int,
	stats WinScoreStats,
) scoreDeltaProbDist {
	pointFreqs := stats.NonDealerWinPointFreqs()
	if actorID == dealerID {
		pointFreqs = stats.DealerWinPointFreqs()
	}
	return multiplyScalarScoreDeltaProbDists(
		winPointsDistFromValidatedStats(pointFreqs),
		winScoreFactorDist(actorID, dealerID, float64(stats.NumSelfDrawWins())/float64(stats.NumWins())),
	)
}

func winScoreDeltaDistFromPointsDist(
	actorID int,
	dealerID int,
	stats WinScoreStats,
	pointsDist scalarProbDist,
) scoreDeltaProbDist {
	return multiplyScalarScoreDeltaProbDists(
		pointsDist,
		winScoreFactorDist(actorID, dealerID, float64(stats.NumSelfDrawWins())/float64(stats.NumWins())),
	)
}
