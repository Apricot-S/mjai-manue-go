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
