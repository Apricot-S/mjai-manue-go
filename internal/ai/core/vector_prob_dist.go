package core

import "slices"

type VectorProbDist struct {
	dist HashMap[[4]float64]
}

func NewVectorProbDist(arg HashMap[[4]float64]) *VectorProbDist {
	pd := &VectorProbDist{dist: NewHashMap[[4]float64]()}
	arg.ForEach(func(value [4]float64, prob float64) {
		if prob > 0 {
			pd.dist.Set(value, prob)
		}
	})
	return pd
}

func (pd *VectorProbDist) Dist() *HashMap[[4]float64] {
	return &pd.dist
}

func (pd *VectorProbDist) Expected() [4]float64 {
	result := [4]float64{0.0, 0.0, 0.0, 0.0}
	pd.dist.ForEach(func(value [4]float64, prob float64) {
		for i, v := range value {
			result[i] += prob * v
		}
	})
	return result
}

func (pd *VectorProbDist) Replace(oldValue [4]float64, newPb *VectorProbDist) *VectorProbDist {
	dist := NewHashMap[[4]float64]()
	prob := 0.0

	pd.dist.ForEach(func(v [4]float64, p float64) {
		if slices.Compare(v[:], oldValue[:]) == 0 {
			prob = p
		} else {
			dist.Set(v, p)
		}
	})

	newPb.dist.ForEach(func(v [4]float64, p float64) {
		dist.Set(v, dist.Get(v, 0.0)+p*prob)
	})

	return &VectorProbDist{dist: dist}
}

func (pd *VectorProbDist) MapValueScalar(mapper func([4]float64) float64) *ScalarProbDist {
	dist := NewHashMap[float64]()
	pd.dist.ForEach(func(v [4]float64, p float64) {
		newValue := mapper(v)
		dist.Set(newValue, dist.Get(newValue, 0.0)+p)
	})
	return &ScalarProbDist{dist: dist}
}

func (pd *VectorProbDist) MapValueVector(mapper func([4]float64) [4]float64) *VectorProbDist {
	dist := NewHashMap[[4]float64]()
	pd.dist.ForEach(func(v [4]float64, p float64) {
		newValue := mapper(v)
		dist.Set(newValue, dist.Get(newValue, 0.0)+p)
	})
	return &VectorProbDist{dist: dist}
}
