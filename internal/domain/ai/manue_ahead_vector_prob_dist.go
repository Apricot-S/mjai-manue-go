package ai

// aheadVector is a four-player 0/1 vector for rank estimation.
//
// An element is 1 when self is estimated to finish ahead of that player, and 0
// otherwise. It is separate from scoreDelta because it represents pairwise rank
// wins, not score changes.
type aheadVector [4]int

type aheadVectorProbDist map[aheadVector]float64

func newAheadVectorProbDist(items map[aheadVector]float64) aheadVectorProbDist {
	dist := make(aheadVectorProbDist, len(items))
	for value, prob := range items {
		if prob > 0.0 {
			dist[value] = prob
		}
	}
	return dist
}

// mapValueScalar maps ahead-vector outcomes to scalar outcomes
// while preserving their probabilities. Outcomes with the same mapped value are
// merged.
func (d aheadVectorProbDist) mapValueScalar(mapper func(aheadVector) float64) scalarProbDist {
	scalars := make(scalarProbDist, len(d))
	for value, prob := range d {
		scalars[mapper(value)] += prob
	}
	return newScalarProbDist(scalars)
}

// mapValueScoreDelta maps ahead-vector outcomes to score-delta outcomes while
// preserving their probabilities. Outcomes with the same mapped value are
// merged.
func (d aheadVectorProbDist) mapValueScoreDelta(mapper func(aheadVector) scoreDelta) scoreDeltaProbDist {
	scoreDeltas := make(scoreDeltaProbDist, len(d))
	for value, prob := range d {
		scoreDeltas[mapper(value)] += prob
	}
	return newScoreDeltaProbDist(scoreDeltas)
}

// addAheadVectorProbDists returns the distribution of lhs + rhs, assuming the
// two ahead-vector random variables are independent.
func addAheadVectorProbDists(lhs, rhs aheadVectorProbDist) aheadVectorProbDist {
	dist := make(aheadVectorProbDist, len(lhs)*len(rhs))
	for lhsValue, lhsProb := range lhs {
		for rhsValue, rhsProb := range rhs {
			var value aheadVector
			for i := range value {
				value[i] = lhsValue[i] + rhsValue[i]
			}
			dist[value] += lhsProb * rhsProb
		}
	}
	return newAheadVectorProbDist(dist)
}

// countAheadWins counts how many opponents self is estimated to finish ahead of.
func countAheadWins(value aheadVector) int {
	sum := 0
	for _, v := range value {
		sum += v
	}
	return sum
}
