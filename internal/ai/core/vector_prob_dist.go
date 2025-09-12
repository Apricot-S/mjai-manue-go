package core

import "slices"

type VectorProbDist struct {
	dist HashMap[[]float64]
}

func NewVectorProbDist(arg HashMap[[]float64]) *VectorProbDist {
	pd := &VectorProbDist{dist: NewHashMap[[]float64]()}
	arg.ForEach(func(value []float64, prob float64) {
		if prob > 0 {
			pd.dist.Set(value, prob)
		}
	})
	return pd
}

func (pd *VectorProbDist) Dist() *HashMap[[]float64] {
	return &pd.dist
}

func (pd *VectorProbDist) Expected() []float64 {
	// result := [4]float64{0.0, 0.0, 0.0, 0.0}
	result := make([]float64, 4) // Assuming 4-dimensional vectors
	pd.dist.ForEach(func(value []float64, prob float64) {
		for i, v := range value {
			result[i] += prob * v
		}
	})
	return result
}

func (pd *VectorProbDist) Replace(oldValue []float64, newPb *VectorProbDist) *VectorProbDist {
	dist := NewHashMap[[]float64]()
	prob := 0.0

	pd.dist.ForEach(func(v []float64, p float64) {
		if slices.Compare(v, oldValue) == 0 {
			prob = p
		} else {
			dist.Set(v, p)
		}
	})

	newPb.dist.ForEach(func(v []float64, p float64) {
		dist.Set(v, dist.Get(v, 0.0)+p*prob)
	})

	return &VectorProbDist{dist: dist}
}

func (pd *VectorProbDist) MapValueScalar(mapper func([]float64) float64) *ScalarProbDist {
	dist := NewHashMap[float64]()
	pd.dist.ForEach(func(v []float64, p float64) {
		newValue := mapper(v)
		dist.Set(newValue, dist.Get(newValue, 0.0)+p)
	})
	return &ScalarProbDist{dist: dist}
}

func (pd *VectorProbDist) MapValueVector(mapper func([]float64) []float64) *VectorProbDist {
	dist := NewHashMap[[]float64]()
	pd.dist.ForEach(func(v []float64, p float64) {
		newValue := mapper(v)
		dist.Set(newValue, dist.Get(newValue, 0.0)+p)
	})
	return &VectorProbDist{dist: dist}
}
