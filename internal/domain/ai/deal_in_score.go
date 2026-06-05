package ai

import "fmt"

type dealInEstimate struct {
	winnerID int
	prob     float64
}

func safeProb(estimates []dealInEstimate) (float64, error) {
	safeProb := 1.0
	for _, estimate := range estimates {
		if estimate.prob < 0.0 || estimate.prob > 1.0 {
			return 0, fmt.Errorf("cannot estimate safe probability: deal-in probability must be between 0 and 1")
		}
		safeProb *= 1.0 - estimate.prob
	}
	return safeProb, nil
}

// immediateDealInScoreDeltaDist returns the immediate score-change
// distribution for a discard against one possible ron winner.
//
// dealInProb is the probability of dealing in to winnerID. The no-change branch
// means this opponent did not ron the discard and the round can continue.
func immediateDealInScoreDeltaDist(
	winnerID int,
	selfID int,
	dealInProb float64,
	pointsDist scalarProbDist,
) (scoreDeltaProbDist, error) {
	if dealInProb < 0.0 || dealInProb > 1.0 {
		return nil, fmt.Errorf("cannot build immediate deal-in score delta distribution: deal-in probability must be between 0 and 1")
	}

	var dealInFactor scoreDelta
	dealInFactor[winnerID] = 1.0
	dealInFactor[selfID] = -1.0
	unitDist := newScoreDeltaProbDist(map[scoreDelta]float64{
		dealInFactor: dealInProb,
		{}:           1.0 - dealInProb,
	})
	return multiplyScalarScoreDeltaProbDists(pointsDist, unitDist), nil
}

func immediateDealInScoreDeltaDistFromStats(
	winnerID int,
	selfID int,
	dealerID int,
	dealInProb float64,
	stats WinScoreStats,
) (scoreDeltaProbDist, error) {
	pointFreqs := stats.NonDealerWinPointFreqs()
	if winnerID == dealerID {
		pointFreqs = stats.DealerWinPointFreqs()
	}
	pointsDist := winPointsDist(pointFreqs)
	return immediateDealInScoreDeltaDist(winnerID, selfID, dealInProb, pointsDist)
}

func immediateScoreDeltaDistFromStats(
	selfID int,
	dealerID int,
	estimates []dealInEstimate,
	stats WinScoreStats,
) (scoreDeltaProbDist, error) {
	dists := make([]scoreDeltaProbDist, 0, len(estimates))
	for _, estimate := range estimates {
		dist, err := immediateDealInScoreDeltaDistFromStats(
			estimate.winnerID,
			selfID,
			dealerID,
			estimate.prob,
			stats,
		)
		if err != nil {
			return nil, err
		}
		dists = append(dists, dist)
	}
	return immediateScoreDeltaDist(dists), nil
}

// immediateScoreDeltaDist merges immediate deal-in distributions from multiple
// possible ron winners.
//
// Each distribution must include a no-change branch. The merge replaces only
// that branch, so it models the same first-ron approximation as Manue and avoids
// expanding double/triple ron combinations.
func immediateScoreDeltaDist(dists []scoreDeltaProbDist) scoreDeltaProbDist {
	result := scoreDeltaProbDist{{}: 1.0}
	for _, dist := range dists {
		result = result.replace(scoreDelta{}, dist)
	}
	return result
}
