package ai

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"

// scoreDelta is a NumPlayers-length score change vector.
type scoreDelta [common.NumPlayers]float64

// scoreDeltaProbDist is a probability distribution of score change vectors.
//
// In Manue's evaluation, these distributions represent how all players' scores
// may change after a candidate action.
type scoreDeltaProbDist map[scoreDelta]float64

type weightedScoreDeltaProbDist struct {
	dist scoreDeltaProbDist
	prob float64
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
