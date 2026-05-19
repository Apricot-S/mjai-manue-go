package ai

type scalarProbDist map[float64]float64

// scoreDelta is a four-player score change vector.
type scoreDelta [4]float64

// scoreDeltaProbDist is a probability distribution of score change vectors.
//
// In Manue's evaluation, these distributions represent how all players' scores
// may change after a candidate action.
type scoreDeltaProbDist map[scoreDelta]float64

type weightedScoreDeltaProbDist struct {
	dist scoreDeltaProbDist
	prob float64
}

// aheadVector is a four-player 0/1 vector for rank estimation.
//
// An element is 1 when self is estimated to finish ahead of that player, and 0
// otherwise. It is separate from scoreDelta because it represents pairwise rank
// wins, not score changes.
type aheadVector [4]int

type aheadVectorProbDist map[aheadVector]float64

// newScalarProbDist builds a scalar probability distribution and drops
// non-positive probabilities.
func newScalarProbDist(items map[float64]float64) scalarProbDist {
	dist := make(scalarProbDist, len(items))
	for value, prob := range items {
		if prob > 0.0 {
			dist[value] = prob
		}
	}
	return dist
}

// expected returns the expected scalar value of the distribution.
func (d scalarProbDist) expected() float64 {
	result := 0.0
	for value, prob := range d {
		result += prob * value
	}
	return result
}

// newScoreDeltaProbDist builds a score-delta probability distribution and drops
// non-positive probabilities.
func newScoreDeltaProbDist(items map[scoreDelta]float64) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(items))
	for value, prob := range items {
		if prob > 0.0 {
			dist[value] = prob
		}
	}
	return dist
}

// expected returns the expected score change vector.
//
// It is sum(probability * scoreDelta) element by element.
func (d scoreDeltaProbDist) expected() scoreDelta {
	var result scoreDelta
	for value, prob := range d {
		for i, v := range value {
			result[i] += prob * v
		}
	}
	return result
}

// replace expands one outcome into another distribution.
//
// In Manue this connects immediate and future score changes. For example,
// immediateScoreChangesDist contains either deal-in score changes or noChanges.
// The noChanges branch means "the discard did not deal in, so the round
// continues"; replace(noChanges, futureScoreChangesDist) redistributes that
// branch's probability mass across the future end-of-round outcomes.
func (d scoreDeltaProbDist) replace(oldValue scoreDelta, newDist scoreDeltaProbDist) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(d)+len(newDist))
	prob := 0.0

	for value, p := range d {
		if value == oldValue {
			prob = p
			continue
		}
		dist[value] = p
	}

	// newDist is conditional on oldValue having happened, so each new outcome
	// gets multiplied by P(oldValue) before returning to the total distribution.
	for value, p := range newDist {
		dist[value] += p * prob
	}
	return newScoreDeltaProbDist(dist)
}

// mapValueScalar maps score-delta outcomes to scalar outcomes while preserving
// their probabilities. Outcomes that map to the same scalar value are merged.
func (d scoreDeltaProbDist) mapValueScalar(mapper func(scoreDelta) float64) scalarProbDist {
	dist := make(scalarProbDist, len(d))
	for value, prob := range d {
		dist[mapper(value)] += prob
	}
	return newScalarProbDist(dist)
}

// mapValueScoreDelta maps score-delta outcomes to other score-delta outcomes
// while preserving their probabilities. Outcomes with the same mapped value are
// merged.
func (d scoreDeltaProbDist) mapValueScoreDelta(mapper func(scoreDelta) scoreDelta) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(d))
	for value, prob := range d {
		dist[mapper(value)] += prob
	}
	return newScoreDeltaProbDist(dist)
}

// addScoreDeltaProbDists returns the distribution of lhs + rhs, assuming the two
// score-delta random variables are independent.
func addScoreDeltaProbDists(lhs, rhs scoreDeltaProbDist) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(lhs)*len(rhs))
	for lhsValue, lhsProb := range lhs {
		for rhsValue, rhsProb := range rhs {
			var value scoreDelta
			for i := range value {
				value[i] = lhsValue[i] + rhsValue[i]
			}
			dist[value] += lhsProb * rhsProb
		}
	}
	return newScoreDeltaProbDist(dist)
}

// multiplyScalarScoreDeltaProbDists returns the distribution of scalar * vector,
// assuming the scalar and score-delta random variables are independent.
func multiplyScalarScoreDeltaProbDists(lhs scalarProbDist, rhs scoreDeltaProbDist) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(lhs)*len(rhs))
	for lhsValue, lhsProb := range lhs {
		for rhsValue, rhsProb := range rhs {
			var value scoreDelta
			for i := range value {
				value[i] = lhsValue * rhsValue[i]
			}
			dist[value] += lhsProb * rhsProb
		}
	}
	return newScoreDeltaProbDist(dist)
}

// mergeScoreDeltaProbDists mixes conditional distributions.
//
// Each item says "use this distribution with this probability". The result is
// the unconditional distribution after weighting and summing all items.
func mergeScoreDeltaProbDists(items []weightedScoreDeltaProbDist) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist)
	for _, item := range items {
		for value, prob := range item.dist {
			dist[value] += prob * item.prob
		}
	}
	return newScoreDeltaProbDist(dist)
}

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
