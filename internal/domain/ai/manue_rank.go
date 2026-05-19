package ai

import "strconv"

type relativeWinProbTable map[string]float64

type rankOpponent struct {
	id       int
	score    float64
	position int
	winProbs relativeWinProbTable
}

// averageRank returns self's expected final rank from pairwise win probabilities
// against the other players.
func averageRank(
	scoreChanges scoreDeltaProbDist,
	selfID int,
	selfScore float64,
	selfPosition int,
	opponents []rankOpponent,
) float64 {
	winsDist := aheadVectorProbDist{{}: 1.0}
	for _, opponent := range opponents {
		winProb := winProbAgainst(
			scoreChanges,
			selfID,
			opponent.id,
			selfScore,
			opponent.score,
			selfPosition,
			opponent.position,
			opponent.winProbs,
		)

		var wins aheadVector
		wins[opponent.id] = 1
		winsDist = addAheadVectorProbDists(winsDist, newAheadVectorProbDist(map[aheadVector]float64{
			{}:   1.0 - winProb,
			wins: winProb,
		}))
	}

	rankDist := winsDist.mapValueScalar(func(wins aheadVector) float64 {
		return float64(4 - countAheadWins(wins))
	})
	return rankDist.expected()
}

// winProbAgainst returns the probability that self finishes ahead of another
// player after applying a score-delta distribution.
func winProbAgainst(
	scoreChanges scoreDeltaProbDist,
	selfID int,
	otherID int,
	selfScore float64,
	otherScore float64,
	selfPosition int,
	otherPosition int,
	winProbs relativeWinProbTable,
) float64 {
	relativeScoreDist := scoreChanges.mapValueScalar(func(scoreChange scoreDelta) float64 {
		return (selfScore + scoreChange[selfID]) - (otherScore + scoreChange[otherID])
	})

	winProb := 0.0
	for relativeScore, prob := range relativeScoreDist {
		winProb += prob * winProbFromRelativeScore(
			relativeScore,
			winProbs,
			selfPosition,
			otherPosition,
		)
	}
	return winProb
}

// winProbFromRelativeScore returns the probability that the player finishes
// ahead of another player from their relative score.
//
// When statistical data is missing, ties are broken by seating order from the
// starting dealer. The player closer to the starting dealer wins a tie.
func winProbFromRelativeScore(
	relativeScore float64,
	winProbs relativeWinProbTable,
	selfPosition int,
	otherPosition int,
) float64 {
	if winProbs != nil {
		if prob, ok := winProbs[strconv.FormatFloat(relativeScore, 'f', 0, 64)]; ok {
			return prob
		}
	}

	if selfPosition < otherPosition {
		if relativeScore >= 0.0 {
			return 1.0
		}
		return 0.0
	}
	if relativeScore > 0.0 {
		return 1.0
	}
	return 0.0
}
