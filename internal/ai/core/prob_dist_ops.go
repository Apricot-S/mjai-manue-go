package core

// Returns probability distribution of sum of two random variables assuming these two are independent.
func AddVectorVector(lhs *VectorProbDist, rhs *VectorProbDist) *VectorProbDist {
	dist := NewHashMap[[]float64]()
	lhs.dist.ForEach(func(v1 []float64, p1 float64) {
		rhs.dist.ForEach(func(v2 []float64, p2 float64) {
			v := make([]float64, 4) // Assuming 4-dimensional vectors
			for i := range len(v) {
				v[i] = v1[i] + v2[i]
			}
			dist.Set(v, dist.Get(v, 0)+p1*p2)
		})
	})
	return &VectorProbDist{dist: dist}
}

// Returns probability distribution of product of two random variables assuming these two are independent.
func MultScalarVector(lhs *ScalarProbDist, rhs *VectorProbDist) *VectorProbDist {
	dist := NewHashMap[[]float64]()
	lhs.dist.ForEach(func(v1 float64, p1 float64) {
		rhs.dist.ForEach(func(v2 []float64, p2 float64) {
			v := make([]float64, 4) // Assuming 4-dimensional vectors
			for i := range len(v) {
				v[i] = v1 * v2[i]
			}
			dist.Set(v, dist.Get(v, 0)+p1*p2)
		})
	})
	return &VectorProbDist{dist: dist}
}

type WeightedVectorProbDist struct {
	Pd   *VectorProbDist
	Prob float64
}

// Merge({{probDist1, prob1}, {probdist2, prob2}, ...})
// Returns a probability distribution of a random variable which follows probDist1 in prob1
// and follows probDist2 in prob2 etc.
func MergeVector(items []WeightedVectorProbDist) *VectorProbDist {
	dist := NewHashMap[[]float64]()
	for _, item := range items {
		item.Pd.dist.ForEach(func(v []float64, p float64) {
			dist.Set(v, dist.Get(v, 0)+p*item.Prob)
		})
	}
	return &VectorProbDist{dist: dist}
}
