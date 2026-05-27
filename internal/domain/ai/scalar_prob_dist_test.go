package ai

import "testing"

func TestScalarProbDist_expected(t *testing.T) {
	dist := newScalarProbDist(map[float64]float64{
		-1000: 0.25,
		2000:  0.75,
		3000:  0,
	})

	if got, want := dist.expected(), 1250.0; got != want {
		t.Errorf("expected() = %v, want %v", got, want)
	}
	if _, ok := dist[3000]; ok {
		t.Errorf("newScalarProbDist() kept zero-probability item")
	}
}
