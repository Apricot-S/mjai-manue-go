package ai

import "testing"

func TestSafeProb(t *testing.T) {
	got, err := safeProb([]dealInEstimate{
		{winnerID: 1, prob: 0.2},
		{winnerID: 2, prob: 0.25},
	})
	if err != nil {
		t.Fatalf("safeProb() failed: %v", err)
	}

	want := 0.6
	if !almostEqual(got, want) {
		t.Errorf("safeProb() = %v, want %v", got, want)
	}
}

func TestSafeProb_ReturnsOneWithoutEstimates(t *testing.T) {
	got, err := safeProb(nil)
	if err != nil {
		t.Fatalf("safeProb() failed: %v", err)
	}
	if got != 1 {
		t.Errorf("safeProb() = %v, want 1", got)
	}
}

func TestSafeProb_ReturnsErrorWithInvalidProb(t *testing.T) {
	_, err := safeProb([]dealInEstimate{{winnerID: 1, prob: -0.1}})
	if err == nil {
		t.Fatal("safeProb() succeeded unexpectedly")
	}
}

func TestImmediateDealInScoreDeltaDist(t *testing.T) {
	got, err := immediateDealInScoreDeltaDist(2, 0, 0.25, scalarProbDist{
		1000: 0.4,
		2000: 0.6,
	})
	if err != nil {
		t.Fatalf("immediateDealInScoreDeltaDist() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{}:                  0.75,
		{-1000, 0, 1000, 0}: 0.10,
		{-2000, 0, 2000, 0}: 0.15,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestImmediateDealInScoreDeltaDist_ReturnsErrorWithInvalidDealInProb(t *testing.T) {
	_, err := immediateDealInScoreDeltaDist(2, 0, -0.1, scalarProbDist{1000: 1})
	if err == nil {
		t.Fatal("immediateDealInScoreDeltaDist() succeeded unexpectedly")
	}
}

func TestImmediateDealInScoreDeltaDistFromStats_SelectsDealerPointFreqs(t *testing.T) {
	got, err := immediateDealInScoreDeltaDistFromStats(2, 0, 2, 0.25, stubManueStats{
		nonDealerWinPointFreqs: map[string]int{
			"1000":  1,
			"total": 1,
		},
		dealerWinPointFreqs: map[string]int{
			"2000":  1,
			"total": 1,
		},
	})
	if err != nil {
		t.Fatalf("immediateDealInScoreDeltaDistFromStats() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{}:                  0.75,
		{-2000, 0, 2000, 0}: 0.25,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestImmediateDealInScoreDeltaDistFromStats_SelectsNonDealerPointFreqs(t *testing.T) {
	got, err := immediateDealInScoreDeltaDistFromStats(1, 0, 2, 0.25, stubManueStats{
		nonDealerWinPointFreqs: map[string]int{
			"1000":  1,
			"total": 1,
		},
		dealerWinPointFreqs: map[string]int{
			"2000":  1,
			"total": 1,
		},
	})
	if err != nil {
		t.Fatalf("immediateDealInScoreDeltaDistFromStats() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{}:               0.75,
		{-1000, 1000, 0}: 0.25,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestImmediateScoreDeltaDist(t *testing.T) {
	got := immediateScoreDeltaDist([]scoreDeltaProbDist{
		{
			{}:            0.8,
			{-1000, 1000}: 0.2,
		},
		{
			{}:               0.75,
			{-2000, 0, 2000}: 0.25,
		},
	})

	want := scoreDeltaProbDist{
		{-1000, 1000}:    0.2,
		{}:               0.6,
		{-2000, 0, 2000}: 0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestImmediateScoreDeltaDistFromStats(t *testing.T) {
	got, err := immediateScoreDeltaDistFromStats(0, 2, []dealInEstimate{
		{winnerID: 1, prob: 0.2},
		{winnerID: 2, prob: 0.25},
	}, stubManueStats{
		nonDealerWinPointFreqs: map[string]int{
			"1000":  1,
			"total": 1,
		},
		dealerWinPointFreqs: map[string]int{
			"2000":  1,
			"total": 1,
		},
	})
	if err != nil {
		t.Fatalf("immediateScoreDeltaDistFromStats() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{-1000, 1000, 0}:    0.2,
		{}:                  0.6,
		{-2000, 0, 2000, 0}: 0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestImmediateScoreDeltaDistFromStats_ReturnsErrorWithInvalidEstimate(t *testing.T) {
	_, err := immediateScoreDeltaDistFromStats(0, 2, []dealInEstimate{
		{winnerID: 1, prob: 1.1},
	}, stubManueStats{
		nonDealerWinPointFreqs: map[string]int{
			"1000":  1,
			"total": 1,
		},
	})
	if err == nil {
		t.Fatal("immediateScoreDeltaDistFromStats() succeeded unexpectedly")
	}
}

func TestImmediateScoreDeltaDist_ReturnsNoChangeWithoutDealInDists(t *testing.T) {
	got := immediateScoreDeltaDist(nil)
	want := scoreDeltaProbDist{{}: 1}
	assertScoreDeltaProbDist(t, got, want)
}
