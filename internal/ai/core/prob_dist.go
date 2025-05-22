package core

import "slices"

type ProbDist[T HashMapKey] struct {
	dist HashMap[T]
}

func NewProbDist[T HashMapKey](arg HashMap[T]) *ProbDist[T] {
	pd := &ProbDist[T]{
		dist: NewHashMap[T](),
	}
	arg.ForEach(func(value T, prob float64) {
		if prob > 0 {
			pd.dist.Set(value, prob)
		}
	})
	return pd
}

func (p *ProbDist[T]) Expected() T {
	var result T
	first := true
	p.dist.ForEach(func(value T, prob float64) {
		weighted := multValue[T](prob, value)
		if first {
			result = weighted
			first = false
		} else {
			result = addValue[T](result, weighted)
		}
	})
	return result
}

func (p *ProbDist[T]) Replace(oldValue T, newPb *ProbDist[T]) *ProbDist[T] {
	dist := NewHashMap[T]()
	var prob float64

	p.dist.ForEach(func(v T, p float64) {
		if equal(v, oldValue) {
			prob = p
		} else {
			dist.Set(v, p)
		}
	})

	newPb.dist.ForEach(func(v T, p float64) {
		dist.Set(v, dist.Get(v, 0)+p*prob)
	})

	return &ProbDist[T]{dist: dist}
}

func (p *ProbDist[T]) MapValue(mapper func(T) T) *ProbDist[T] {
	dist := NewHashMap[T]()
	p.dist.ForEach(func(v T, p float64) {
		newValue := mapper(v)
		dist.Set(newValue, dist.Get(newValue, 0.0)+p)
	})
	return &ProbDist[T]{dist: dist}
}

// Returns probability distribution of sum of two random variables assuming these two are independent.
func Add[T, U, V HashMapKey](lhs *ProbDist[T], rhs *ProbDist[U]) *ProbDist[V] {
	dist := NewHashMap[V]()
	lhs.dist.ForEach(func(v1 T, p1 float64) {
		rhs.dist.ForEach(func(v2 U, p2 float64) {
			v := addValue[V](v1, v2)
			dist.Set(v, dist.Get(v, 0)+p1*p2)
		})
	})
	return &ProbDist[V]{dist: dist}
}

// Returns probability distribution of product of two random variables assuming these two are independent.
func Mult[T, U, V HashMapKey](lhs *ProbDist[T], rhs *ProbDist[U]) *ProbDist[V] {
	dist := NewHashMap[V]()
	lhs.dist.ForEach(func(v1 T, p1 float64) {
		rhs.dist.ForEach(func(v2 U, p2 float64) {
			v := multValue[V](v1, v2)
			dist.Set(v, dist.Get(v, 0)+p1*p2)
		})
	})
	return &ProbDist[V]{dist: dist}
}

type WeightedProbDist[T HashMapKey] struct {
	Pd   *ProbDist[T]
	Prob float64
}

// Merge({{probDist1, prob1}, {probdist2, prob2}, ...})
// Returns a probability distribution of a random variable which follows probDist1 in prob1
// and follows probDist2 in prob2 etc.
func Merge[T HashMapKey](items []WeightedProbDist[T]) *ProbDist[T] {
	dist := NewHashMap[T]()
	for _, item := range items {
		item.Pd.dist.ForEach(func(v T, p float64) {
			dist.Set(v, dist.Get(v, 0)+p*item.Prob)
		})
	}
	return &ProbDist[T]{dist: dist}
}

func equal[T float64 | []float64](lhs, rhs T) bool {
	switch l := any(lhs).(type) {
	case float64:
		return l == any(rhs).(float64)
	case []float64:
		return slices.Compare(l, any(rhs).([]float64)) == 0
	default:
		panic("unsupported types for equal")
	}
}

func addValue[T float64 | []float64](lhs, rhs any) T {
	switch l := lhs.(type) {
	case float64:
		switch r := rhs.(type) {
		case float64:
			return any(l + r).(T)
		case []float64:
			result := make([]float64, len(r))
			for i := range r {
				result[i] = l + r[i]
			}
			return any(result).(T)
		default:
			panic("unsupported types for addValue")
		}
	case []float64:
		switch r := rhs.(type) {
		case float64:
			result := make([]float64, len(l))
			for i := range l {
				result[i] = l[i] + r
			}
			return any(result).(T)
		case []float64:
			if len(l) != len(r) {
				panic("length mismatch")
			}

			result := make([]float64, len(l))
			for i := range l {
				result[i] = l[i] + r[i]
			}
			return any(result).(T)
		default:
			panic("unsupported types for addValue")
		}
	default:
		panic("unsupported types for addValue")
	}
}

func multValue[T float64 | []float64](lhs, rhs any) T {
	switch l := lhs.(type) {
	case float64:
		switch r := rhs.(type) {
		case float64:
			return any(l * r).(T)
		case []float64:
			result := make([]float64, len(r))
			for i := range r {
				result[i] = l * r[i]
			}
			return any(result).(T)
		default:
			panic("unsupported types for multValue")
		}
	case []float64:
		switch r := rhs.(type) {
		case float64:
			result := make([]float64, len(l))
			for i := range l {
				result[i] = l[i] * r
			}
			return any(result).(T)
		case []float64:
			if len(l) != len(r) {
				panic("length mismatch")
			}

			result := make([]float64, len(l))
			for i := range l {
				result[i] = l[i] * r[i]
			}
			return any(result).(T)
		default:
			panic("unsupported types for multValue")
		}
	default:
		panic("unsupported types for multValue")
	}
}
