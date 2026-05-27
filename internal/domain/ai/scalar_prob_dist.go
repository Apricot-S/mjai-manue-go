package ai

type scalarProbDist map[float64]float64

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
