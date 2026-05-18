package ai

type scalarProbDist map[float64]float64

type scoreDelta [4]float64

type scoreDeltaProbDist map[scoreDelta]float64

type weightedScoreDeltaProbDist struct {
	dist scoreDeltaProbDist
	prob float64
}

func newScalarProbDist(items map[float64]float64) scalarProbDist {
	dist := make(scalarProbDist, len(items))
	for value, prob := range items {
		if prob > 0.0 {
			dist[value] = prob
		}
	}
	return dist
}

func (d scalarProbDist) expected() float64 {
	result := 0.0
	for value, prob := range d {
		result += prob * value
	}
	return result
}

func newScoreDeltaProbDist(items map[scoreDelta]float64) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(items))
	for value, prob := range items {
		if prob > 0.0 {
			dist[value] = prob
		}
	}
	return dist
}

func (d scoreDeltaProbDist) expected() scoreDelta {
	var result scoreDelta
	for value, prob := range d {
		for i, v := range value {
			result[i] += prob * v
		}
	}
	return result
}

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

	for value, p := range newDist {
		dist[value] += p * prob
	}
	return newScoreDeltaProbDist(dist)
}

func (d scoreDeltaProbDist) mapValueScalar(mapper func(scoreDelta) float64) scalarProbDist {
	dist := make(scalarProbDist, len(d))
	for value, prob := range d {
		dist[mapper(value)] += prob
	}
	return newScalarProbDist(dist)
}

func (d scoreDeltaProbDist) mapValueScoreDelta(mapper func(scoreDelta) scoreDelta) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, len(d))
	for value, prob := range d {
		dist[mapper(value)] += prob
	}
	return newScoreDeltaProbDist(dist)
}

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

func mergeScoreDeltaProbDists(items []weightedScoreDeltaProbDist) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist)
	for _, item := range items {
		for value, prob := range item.dist {
			dist[value] += prob * item.prob
		}
	}
	return newScoreDeltaProbDist(dist)
}
