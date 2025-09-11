package core

type ScalarProbDist struct {
	dist HashMap[float64]
}

func NewScalarProbDist(arg HashMap[float64]) *ScalarProbDist {
	pd := &ScalarProbDist{dist: NewHashMap[float64]()}
	arg.ForEach(func(value float64, prob float64) {
		if prob > 0 {
			pd.dist.Set(value, prob)
		}
	})
	return pd
}

func (p *ScalarProbDist) Dist() *HashMap[float64] {
	return &p.dist
}

func (p *ScalarProbDist) Expected() float64 {
	result := 0.0
	p.dist.ForEach(func(value float64, prob float64) {
		result += prob * value
	})
	return result
}
