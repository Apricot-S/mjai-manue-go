package ai

import (
	"math"
	"testing"
)

func assertScalarProbDist(t *testing.T, got, want scalarProbDist) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("len(dist) = %d, want %d; dist = %v", len(got), len(want), got)
	}
	for value, wantProb := range want {
		if gotProb := got[value]; !almostEqual(gotProb, wantProb) {
			t.Errorf("dist[%v] = %v, want %v", value, gotProb, wantProb)
		}
	}
}

func assertScoreDeltaProbDist(t *testing.T, got, want scoreDeltaProbDist) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("len(dist) = %d, want %d; dist = %v", len(got), len(want), got)
	}
	for value, wantProb := range want {
		if gotProb := got[value]; !almostEqual(gotProb, wantProb) {
			t.Errorf("dist[%v] = %v, want %v", value, gotProb, wantProb)
		}
	}
}

func almostEqual(lhs, rhs float64) bool {
	const epsilon = 1e-12
	return math.Abs(lhs-rhs) <= epsilon
}
