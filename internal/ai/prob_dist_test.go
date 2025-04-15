package ai

import (
	"math"
	"testing"
)

func assertFloatEqual(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) >= 0.001 {
		t.Errorf("Expected %f, but got %f", want, got)
	}
}

func TestProbDist(t *testing.T) {
	hm1 := NewHashMap[float64]()
	hm1.Set(0, 0.5)
	hm1.Set(8000, 0.5)
	pb1 := NewProbDist(hm1)

	hm2 := NewHashMap[float64]()
	hm2.Set(0, 0.5)
	hm2.Set(-2000, 0.5)
	pb2 := NewProbDist(hm2)

	hm3 := NewHashMap[[]float64]()
	hm3.Set([]float64{0, 1}, 0.5)
	hm3.Set([]float64{2, 3}, 0.5)
	pb3 := NewProbDist(hm3)

	hm4 := NewHashMap[[]float64]()
	hm4.Set([]float64{1, 2}, 0.5)
	hm4.Set([]float64{4, 8}, 0.5)
	pb4 := NewProbDist(hm4)

	hm5 := NewHashMap[float64]()
	hm5.Set(1, 0.5)
	hm5.Set(-1, 0.5)
	pb5 := NewProbDist(hm5)

	t.Run("Add single distributions", func(t *testing.T) {
		rpb1 := Add[float64, float64, float64](pb1, pb1)

		assertFloatEqual(t, rpb1.dist.Get(0, 0), 0.25)
		assertFloatEqual(t, rpb1.dist.Get(8000, 0), 0.5)
		assertFloatEqual(t, rpb1.dist.Get(16000, 0), 0.25)
		assertFloatEqual(t, rpb1.Expected(), 8000)
	})

	t.Run("Merge distributions", func(t *testing.T) {
		rpb2 := Merge([]WeightedProbDist[float64]{{pb1, 0.5}, {pb2, 0.5}})

		assertFloatEqual(t, rpb2.dist.Get(0, 0), 0.5)
		assertFloatEqual(t, rpb2.dist.Get(8000, 0), 0.25)
		assertFloatEqual(t, rpb2.dist.Get(-2000, 0), 0.25)
		assertFloatEqual(t, rpb2.Expected(), 8000*0.25-2000*0.25)
	})

	t.Run("Add array distributions", func(t *testing.T) {
		rpb3 := Add[[]float64, []float64, []float64](pb3, pb3)

		assertFloatEqual(t, rpb3.dist.Get([]float64{0, 2}, 0), 0.25)
		assertFloatEqual(t, rpb3.dist.Get([]float64{2, 4}, 0), 0.5)
		assertFloatEqual(t, rpb3.dist.Get([]float64{4, 6}, 0), 0.25)

		expected := rpb3.Expected()
		assertFloatEqual(t, expected[0], 2)
		assertFloatEqual(t, expected[1], 4)
	})

	t.Run("Multiply distributions", func(t *testing.T) {
		rpb4 := Mult[[]float64, float64, []float64](pb4, pb5)

		assertFloatEqual(t, rpb4.dist.Get([]float64{1, 2}, 0), 0.25)
		assertFloatEqual(t, rpb4.dist.Get([]float64{-1, -2}, 0), 0.25)
		assertFloatEqual(t, rpb4.dist.Get([]float64{4, 8}, 0), 0.25)
		assertFloatEqual(t, rpb4.dist.Get([]float64{-4, -8}, 0), 0.25)
	})
}
