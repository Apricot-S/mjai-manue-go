package ai

import "testing"

func TestFutureScoreDeltaDist(t *testing.T) {
	selfWinDist := scoreDeltaProbDist{{1000, 0, 0, 0}: 1}
	exhaustiveDrawDist := scoreDeltaProbDist{{0, 1000, 0, 0}: 1}
	otherWinDists := []scoreDeltaProbDist{
		{{0, 0, 1000, 0}: 1},
		{{0, 0, 0, 1000}: 1},
	}

	got := futureScoreDeltaDist(selfWinDist, 0.2, exhaustiveDrawDist, 0.3, otherWinDists, 0.5)
	want := scoreDeltaProbDist{
		{1000, 0, 0, 0}: 0.2,
		{0, 1000, 0, 0}: 0.3,
		{0, 0, 1000, 0}: 0.25,
		{0, 0, 0, 1000}: 0.25,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestTotalScoreDeltaDist(t *testing.T) {
	immediateDist := scoreDeltaProbDist{
		{}:            0.8,
		{-1000, 1000}: 0.2,
	}
	futureDist := scoreDeltaProbDist{
		{1000, 0, 0, 0}: 0.25,
		{0, 1000, 0, 0}: 0.75,
	}

	got := totalScoreDeltaDist(immediateDist, futureDist)
	want := scoreDeltaProbDist{
		{-1000, 1000}:   0.2,
		{1000, 0, 0, 0}: 0.2,
		{0, 1000, 0, 0}: 0.6,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestExpectedPts(t *testing.T) {
	got := expectedPts(0, scoreDeltaProbDist{
		{1000, -1000, 0, 0}: 0.25,
		{-500, 500, 0, 0}:   0.75,
	})
	want := -125.0
	if got != want {
		t.Errorf("expectedPts() = %v, want %v", got, want)
	}
}
