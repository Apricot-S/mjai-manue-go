package ai

import "testing"

func TestScoreDeltaProbDist_expected(t *testing.T) {
	dist := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}:     0.25,
		{2000, 0, -1000, -1000}: 0.75,
	})

	want := scoreDelta{1750, -250, -750, -750}
	if got := dist.expected(); got != want {
		t.Errorf("expected() = %v, want %v", got, want)
	}
}

func TestScoreDeltaProbDist_replace(t *testing.T) {
	noChange := scoreDelta{}
	dist := newScoreDeltaProbDist(map[scoreDelta]float64{
		noChange:            0.8,
		{1000, -1000, 0, 0}: 0.2,
	})
	replacement := newScoreDeltaProbDist(map[scoreDelta]float64{
		{2000, 0, -1000, -1000}:     0.25,
		{3000, -1000, -1000, -1000}: 0.75,
	})

	got := dist.replace(noChange, replacement)
	want := scoreDeltaProbDist{
		{1000, -1000, 0, 0}:         0.2,
		{2000, 0, -1000, -1000}:     0.2,
		{3000, -1000, -1000, -1000}: 0.6,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestScoreDeltaProbDist_replace_DropsNonPositiveProbability(t *testing.T) {
	noChange := scoreDelta{}
	dist := newScoreDeltaProbDist(map[scoreDelta]float64{
		noChange:            0.8,
		{1000, -1000, 0, 0}: 0.2,
	})
	replacement := scoreDeltaProbDist{
		{2000, 0, -1000, -1000}:     0.25,
		{3000, -1000, -1000, -1000}: -0.75,
	}

	got := dist.replace(noChange, replacement)
	want := scoreDeltaProbDist{
		{1000, -1000, 0, 0}:     0.2,
		{2000, 0, -1000, -1000}: 0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestMultiplyScalarScoreDeltaProbDists(t *testing.T) {
	lhs := newScalarProbDist(map[float64]float64{2: 0.25, 3: 0.75})
	rhs := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}:     0.4,
		{2000, 0, -1000, -1000}: 0.6,
	})

	got := multiplyScalarScoreDeltaProbDists(lhs, rhs)
	want := scoreDeltaProbDist{
		{2000, -2000, 0, 0}:     0.10,
		{4000, 0, -2000, -2000}: 0.15,
		{3000, -3000, 0, 0}:     0.30,
		{6000, 0, -3000, -3000}: 0.45,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestMergeScoreDeltaProbDists(t *testing.T) {
	got := mergeScoreDeltaProbDists([]weightedScoreDeltaProbDist{
		{
			dist: newScoreDeltaProbDist(map[scoreDelta]float64{
				{1000, -1000, 0, 0}: 1,
			}),
			prob: 0.25,
		},
		{
			dist: newScoreDeltaProbDist(map[scoreDelta]float64{
				{2000, 0, -1000, -1000}:     0.4,
				{3000, -1000, -1000, -1000}: 0.6,
			}),
			prob: 0.75,
		},
	})
	want := scoreDeltaProbDist{
		{1000, -1000, 0, 0}:         0.25,
		{2000, 0, -1000, -1000}:     0.30,
		{3000, -1000, -1000, -1000}: 0.45,
	}
	assertScoreDeltaProbDist(t, got, want)
}
