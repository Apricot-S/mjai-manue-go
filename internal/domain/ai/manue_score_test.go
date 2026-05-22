package ai

import "testing"

func TestWinScoreFactor(t *testing.T) {
	tests := []struct {
		name     string
		actorID  int
		targetID int
		dealerID int
		want     scoreDelta
	}{
		{
			name:     "ron",
			actorID:  1,
			targetID: 2,
			dealerID: 0,
			want:     scoreDelta{0, 1, -1, 0},
		},
		{
			name:     "dealer self draw",
			actorID:  0,
			targetID: 0,
			dealerID: 0,
			want:     scoreDelta{1, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0},
		},
		{
			name:     "non-dealer self draw",
			actorID:  1,
			targetID: 1,
			dealerID: 0,
			want:     scoreDelta{-1.0 / 2.0, 1, -1.0 / 4.0, -1.0 / 4.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := winScoreFactor(tt.actorID, tt.targetID, tt.dealerID)
			if got != tt.want {
				t.Errorf("winScoreFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWinEstimate(t *testing.T) {
	got, err := newWinEstimate(100, 25, 100000, map[float64]int{
		2000: 10,
		5000: 15,
	})
	if err != nil {
		t.Fatalf("newWinEstimate() failed: %v", err)
	}

	if !almostEqual(got.prob, 0.25) {
		t.Errorf("prob = %v, want 0.25", got.prob)
	}
	if !almostEqual(got.avgPts, 4000) {
		t.Errorf("avgPts = %v, want 4000", got.avgPts)
	}
	if !almostEqual(got.expPts, 1000) {
		t.Errorf("expPts = %v, want 1000", got.expPts)
	}
	assertScalarProbDist(t, got.pointsDist, scalarProbDist{
		2000: 0.4,
		5000: 0.6,
	})
}

func TestNewWinEstimate_ReturnsZeroWithoutWins(t *testing.T) {
	got, err := newWinEstimate(100, 0, 0, nil)
	if err != nil {
		t.Fatalf("newWinEstimate() failed: %v", err)
	}

	if got.prob != 0 {
		t.Errorf("prob = %v, want 0", got.prob)
	}
	if got.avgPts != 0 {
		t.Errorf("avgPts = %v, want 0", got.avgPts)
	}
	if got.expPts != 0 {
		t.Errorf("expPts = %v, want 0", got.expPts)
	}
	assertScalarProbDist(t, got.pointsDist, scalarProbDist{})
}

func TestNewWinEstimate_ReturnsErrorWithInvalidInputs(t *testing.T) {
	tests := []struct {
		name       string
		numTries   int
		totalWins  int
		totalPts   float64
		pointFreqs map[float64]int
	}{
		{
			name:      "invalid numTries",
			numTries:  0,
			totalWins: 1,
			totalPts:  1000,
		},
		{
			name:      "negative totalWins",
			numTries:  100,
			totalWins: -1,
		},
		{
			name:      "too many totalWins",
			numTries:  100,
			totalWins: 101,
		},
		{
			name:      "negative totalPts",
			numTries:  100,
			totalWins: 1,
			totalPts:  -1,
		},
		{
			name:       "zero wins with points",
			numTries:   100,
			totalWins:  0,
			totalPts:   1000,
			pointFreqs: map[float64]int{1000: 1},
		},
		{
			name:       "non-positive points",
			numTries:   100,
			totalWins:  1,
			totalPts:   1000,
			pointFreqs: map[float64]int{0: 1},
		},
		{
			name:       "negative frequency",
			numTries:   100,
			totalWins:  1,
			totalPts:   1000,
			pointFreqs: map[float64]int{1000: -1},
		},
		{
			name:       "frequency sum mismatch",
			numTries:   100,
			totalWins:  2,
			totalPts:   2000,
			pointFreqs: map[float64]int{1000: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newWinEstimate(tt.numTries, tt.totalWins, tt.totalPts, tt.pointFreqs)
			if err == nil {
				t.Fatal("newWinEstimate() succeeded unexpectedly")
			}
		})
	}
}

func TestWinEstimateAccumulator(t *testing.T) {
	var accumulator winEstimateAccumulator
	accumulator.addNoWinTrial()
	if err := accumulator.addWinTrial(2000); err != nil {
		t.Fatalf("addWinTrial(2000) failed: %v", err)
	}
	if err := accumulator.addWinTrial(5000); err != nil {
		t.Fatalf("addWinTrial(5000) failed: %v", err)
	}
	if err := accumulator.addWinTrial(5000); err != nil {
		t.Fatalf("addWinTrial(5000) failed: %v", err)
	}

	got, err := accumulator.estimate()
	if err != nil {
		t.Fatalf("estimate() failed: %v", err)
	}

	if !almostEqual(got.prob, 0.75) {
		t.Errorf("prob = %v, want 0.75", got.prob)
	}
	if !almostEqual(got.avgPts, 4000) {
		t.Errorf("avgPts = %v, want 4000", got.avgPts)
	}
	if !almostEqual(got.expPts, 3000) {
		t.Errorf("expPts = %v, want 3000", got.expPts)
	}
	assertScalarProbDist(t, got.pointsDist, scalarProbDist{
		2000: 1.0 / 3.0,
		5000: 2.0 / 3.0,
	})
}

func TestWinEstimateAccumulator_AddWinTrialReturnsErrorWithInvalidPoints(t *testing.T) {
	var accumulator winEstimateAccumulator
	err := accumulator.addWinTrial(0)
	if err == nil {
		t.Fatal("addWinTrial() succeeded unexpectedly")
	}
}

func TestWinEstimateAccumulator_Merge(t *testing.T) {
	var lhs winEstimateAccumulator
	lhs.addNoWinTrial()
	if err := lhs.addWinTrial(2000); err != nil {
		t.Fatalf("lhs.addWinTrial() failed: %v", err)
	}

	var rhs winEstimateAccumulator
	if err := rhs.addWinTrial(5000); err != nil {
		t.Fatalf("rhs.addWinTrial() failed: %v", err)
	}
	rhs.addNoWinTrial()

	if err := lhs.merge(rhs); err != nil {
		t.Fatalf("merge() failed: %v", err)
	}
	got, err := lhs.estimate()
	if err != nil {
		t.Fatalf("estimate() failed: %v", err)
	}

	if !almostEqual(got.prob, 0.5) {
		t.Errorf("prob = %v, want 0.5", got.prob)
	}
	if !almostEqual(got.avgPts, 3500) {
		t.Errorf("avgPts = %v, want 3500", got.avgPts)
	}
	if !almostEqual(got.expPts, 1750) {
		t.Errorf("expPts = %v, want 1750", got.expPts)
	}
	assertScalarProbDist(t, got.pointsDist, scalarProbDist{
		2000: 0.5,
		5000: 0.5,
	})
}

func TestWinEstimateAccumulator_MergeReturnsErrorWithInvalidAccumulator(t *testing.T) {
	tests := []struct {
		name  string
		other winEstimateAccumulator
	}{
		{
			name: "too many wins",
			other: winEstimateAccumulator{
				numTries:  1,
				totalWins: 2,
			},
		},
		{
			name: "frequency sum mismatch",
			other: winEstimateAccumulator{
				numTries:   2,
				totalWins:  2,
				totalPts:   2000,
				pointFreqs: map[float64]int{1000: 1},
			},
		},
		{
			name: "invalid points",
			other: winEstimateAccumulator{
				numTries:   1,
				totalWins:  1,
				totalPts:   1000,
				pointFreqs: map[float64]int{0: 1},
			},
		},
		{
			name: "negative frequency",
			other: winEstimateAccumulator{
				numTries:   1,
				totalWins:  1,
				totalPts:   1000,
				pointFreqs: map[float64]int{1000: -1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accumulator := winEstimateAccumulator{
				numTries:   1,
				totalWins:  1,
				totalPts:   1000,
				pointFreqs: map[float64]int{1000: 1},
			}
			err := accumulator.merge(tt.other)
			if err == nil {
				t.Fatal("merge() succeeded unexpectedly")
			}
			if accumulator.numTries != 1 {
				t.Errorf("numTries = %v, want unchanged 1", accumulator.numTries)
			}
			if accumulator.totalWins != 1 {
				t.Errorf("totalWins = %v, want unchanged 1", accumulator.totalWins)
			}
			if accumulator.totalPts != 1000 {
				t.Errorf("totalPts = %v, want unchanged 1000", accumulator.totalPts)
			}
			if accumulator.pointFreqs[1000] != 1 {
				t.Errorf("pointFreqs[1000] = %v, want unchanged 1", accumulator.pointFreqs[1000])
			}
		})
	}
}

func TestWinEstimateAccumulatorSet(t *testing.T) {
	accumulators := winEstimateAccumulatorSet{}
	accumulators.addNoWinTrial("discard-1m")
	if err := accumulators.addWinTrial("discard-1m", 2000); err != nil {
		t.Fatalf("addWinTrial(discard-1m) failed: %v", err)
	}
	if err := accumulators.addWinTrial("riichi-1m", 5000); err != nil {
		t.Fatalf("addWinTrial(riichi-1m) failed: %v", err)
	}

	got, err := accumulators.estimates()
	if err != nil {
		t.Fatalf("estimates() failed: %v", err)
	}

	discardEstimate := got["discard-1m"]
	if !almostEqual(discardEstimate.prob, 0.5) {
		t.Errorf("discard prob = %v, want 0.5", discardEstimate.prob)
	}
	if !almostEqual(discardEstimate.avgPts, 2000) {
		t.Errorf("discard avgPts = %v, want 2000", discardEstimate.avgPts)
	}
	if !almostEqual(discardEstimate.expPts, 1000) {
		t.Errorf("discard expPts = %v, want 1000", discardEstimate.expPts)
	}

	riichiEstimate := got["riichi-1m"]
	if riichiEstimate.prob != 1 {
		t.Errorf("riichi prob = %v, want 1", riichiEstimate.prob)
	}
	if riichiEstimate.avgPts != 5000 {
		t.Errorf("riichi avgPts = %v, want 5000", riichiEstimate.avgPts)
	}
	if riichiEstimate.expPts != 5000 {
		t.Errorf("riichi expPts = %v, want 5000", riichiEstimate.expPts)
	}
}

func TestNewWinEstimateAccumulatorSet(t *testing.T) {
	got := newWinEstimateAccumulatorSet([]string{"discard-1m", "discard-2m"})

	if len(got) != 2 {
		t.Fatalf("len(newWinEstimateAccumulatorSet()) = %v, want 2", len(got))
	}
	if _, ok := got["discard-1m"]; !ok {
		t.Errorf("newWinEstimateAccumulatorSet() does not contain discard-1m")
	}
	if _, ok := got["discard-2m"]; !ok {
		t.Errorf("newWinEstimateAccumulatorSet() does not contain discard-2m")
	}
}

func TestWinEstimateAccumulatorSet_AddTrial(t *testing.T) {
	accumulators := newWinEstimateAccumulatorSet([]string{"discard-1m", "discard-2m"})

	if err := accumulators.addTrial(map[string]float64{"discard-1m": 2000}); err != nil {
		t.Fatalf("addTrial(first) failed: %v", err)
	}
	if err := accumulators.addTrial(map[string]float64{"discard-2m": 5000}); err != nil {
		t.Fatalf("addTrial(second) failed: %v", err)
	}

	got, err := accumulators.estimates()
	if err != nil {
		t.Fatalf("estimates() failed: %v", err)
	}

	if !almostEqual(got["discard-1m"].prob, 0.5) {
		t.Errorf("discard-1m prob = %v, want 0.5", got["discard-1m"].prob)
	}
	if got["discard-1m"].avgPts != 2000 {
		t.Errorf("discard-1m avgPts = %v, want 2000", got["discard-1m"].avgPts)
	}
	if !almostEqual(got["discard-2m"].prob, 0.5) {
		t.Errorf("discard-2m prob = %v, want 0.5", got["discard-2m"].prob)
	}
	if got["discard-2m"].avgPts != 5000 {
		t.Errorf("discard-2m avgPts = %v, want 5000", got["discard-2m"].avgPts)
	}
}

func TestWinEstimateAccumulatorSet_AddTrialReturnsErrorWithoutPartialUpdate(t *testing.T) {
	tests := []struct {
		name        string
		winPtsByKey map[string]float64
	}{
		{
			name:        "unknown key",
			winPtsByKey: map[string]float64{"unknown": 2000},
		},
		{
			name:        "invalid points",
			winPtsByKey: map[string]float64{"discard-1m": 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accumulators := newWinEstimateAccumulatorSet([]string{"discard-1m"})
			err := accumulators.addTrial(tt.winPtsByKey)
			if err == nil {
				t.Fatal("addTrial() succeeded unexpectedly")
			}
			if accumulators["discard-1m"].numTries != 0 {
				t.Errorf("numTries = %v, want unchanged 0", accumulators["discard-1m"].numTries)
			}
		})
	}
}

func TestWinEstimateAccumulatorSet_Merge(t *testing.T) {
	lhs := winEstimateAccumulatorSet{}
	if err := lhs.addWinTrial("discard-1m", 2000); err != nil {
		t.Fatalf("lhs.addWinTrial() failed: %v", err)
	}
	lhs.addNoWinTrial("discard-2m")

	rhs := winEstimateAccumulatorSet{}
	rhs.addNoWinTrial("discard-1m")
	if err := rhs.addWinTrial("discard-2m", 5000); err != nil {
		t.Fatalf("rhs.addWinTrial() failed: %v", err)
	}

	if err := lhs.merge(rhs); err != nil {
		t.Fatalf("merge() failed: %v", err)
	}
	got, err := lhs.estimates()
	if err != nil {
		t.Fatalf("estimates() failed: %v", err)
	}

	if !almostEqual(got["discard-1m"].prob, 0.5) {
		t.Errorf("discard-1m prob = %v, want 0.5", got["discard-1m"].prob)
	}
	if !almostEqual(got["discard-2m"].prob, 0.5) {
		t.Errorf("discard-2m prob = %v, want 0.5", got["discard-2m"].prob)
	}
}

func TestWinEstimateAccumulatorSet_MergeReturnsErrorWithInvalidAccumulator(t *testing.T) {
	accumulators := winEstimateAccumulatorSet{
		"valid": {
			numTries:   1,
			totalWins:  1,
			totalPts:   1000,
			pointFreqs: map[float64]int{1000: 1},
		},
	}
	err := accumulators.merge(winEstimateAccumulatorSet{
		"valid": {
			numTries: 1,
		},
		"discard-1m": {
			numTries:  1,
			totalWins: 2,
		},
	})
	if err == nil {
		t.Fatal("merge() succeeded unexpectedly")
	}
	if _, ok := accumulators["discard-1m"]; ok {
		t.Fatal("merge() inserted invalid accumulator unexpectedly")
	}
	if accumulators["valid"].numTries != 1 {
		t.Errorf("valid numTries = %v, want unchanged 1", accumulators["valid"].numTries)
	}
	if accumulators["valid"].totalWins != 1 {
		t.Errorf("valid totalWins = %v, want unchanged 1", accumulators["valid"].totalWins)
	}
	if accumulators["valid"].pointFreqs[1000] != 1 {
		t.Errorf("valid pointFreqs[1000] = %v, want unchanged 1", accumulators["valid"].pointFreqs[1000])
	}
}

func TestWinEstimateAccumulatorSet_EstimatesReturnsErrorWithInvalidAccumulator(t *testing.T) {
	accumulators := winEstimateAccumulatorSet{
		"discard-1m": {
			numTries: 0,
		},
	}
	_, err := accumulators.estimates()
	if err == nil {
		t.Fatal("estimates() succeeded unexpectedly")
	}
}

func TestWinEstimatesFromTrials(t *testing.T) {
	got, err := winEstimatesFromTrials(
		[]string{"discard-1m", "discard-2m"},
		[]map[string]float64{
			{"discard-1m": 2000},
			{"discard-1m": 5000, "discard-2m": 8000},
			{},
		},
	)
	if err != nil {
		t.Fatalf("winEstimatesFromTrials() failed: %v", err)
	}

	discard1m := got["discard-1m"]
	if !almostEqual(discard1m.prob, 2.0/3.0) {
		t.Errorf("discard-1m prob = %v, want %v", discard1m.prob, 2.0/3.0)
	}
	if !almostEqual(discard1m.avgPts, 3500) {
		t.Errorf("discard-1m avgPts = %v, want 3500", discard1m.avgPts)
	}
	if !almostEqual(discard1m.expPts, 7000.0/3.0) {
		t.Errorf("discard-1m expPts = %v, want %v", discard1m.expPts, 7000.0/3.0)
	}
	assertScalarProbDist(t, discard1m.pointsDist, scalarProbDist{
		2000: 0.5,
		5000: 0.5,
	})

	discard2m := got["discard-2m"]
	if !almostEqual(discard2m.prob, 1.0/3.0) {
		t.Errorf("discard-2m prob = %v, want %v", discard2m.prob, 1.0/3.0)
	}
	if discard2m.avgPts != 8000 {
		t.Errorf("discard-2m avgPts = %v, want 8000", discard2m.avgPts)
	}
	if !almostEqual(discard2m.expPts, 8000.0/3.0) {
		t.Errorf("discard-2m expPts = %v, want %v", discard2m.expPts, 8000.0/3.0)
	}
}

func TestWinEstimatesFromTrials_ReturnsErrorWithInvalidTrial(t *testing.T) {
	_, err := winEstimatesFromTrials(
		[]string{"discard-1m"},
		[]map[string]float64{{"unknown": 2000}},
	)
	if err == nil {
		t.Fatal("winEstimatesFromTrials() succeeded unexpectedly")
	}
}

func TestDealInProb(t *testing.T) {
	got, err := dealInProb([]dealInEstimate{
		{winnerID: 1, prob: 0.2},
		{winnerID: 2, prob: 0.25},
	})
	if err != nil {
		t.Fatalf("dealInProb() failed: %v", err)
	}

	want := 0.4
	if !almostEqual(got, want) {
		t.Errorf("dealInProb() = %v, want %v", got, want)
	}
}

func TestDealInProb_ReturnsZeroWithoutEstimates(t *testing.T) {
	got, err := dealInProb(nil)
	if err != nil {
		t.Fatalf("dealInProb() failed: %v", err)
	}
	if got != 0 {
		t.Errorf("dealInProb() = %v, want 0", got)
	}
}

func TestDealInProb_ReturnsErrorWithInvalidProb(t *testing.T) {
	_, err := dealInProb([]dealInEstimate{{winnerID: 1, prob: 1.1}})
	if err == nil {
		t.Fatal("dealInProb() succeeded unexpectedly")
	}
}

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

func TestWinScoreFactorDist(t *testing.T) {
	got := winScoreFactorDist(1, 0, 0.4)
	want := scoreDeltaProbDist{
		{-1.0 / 2.0, 1, -1.0 / 4.0, -1.0 / 4.0}: 0.4,
		{-1, 1, 0, 0}:                           0.2,
		{0, 1, -1, 0}:                           0.2,
		{0, 1, 0, -1}:                           0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestWinPointsDist(t *testing.T) {
	got, err := winPointsDist(map[string]int{
		"1000":  25,
		"2000":  75,
		"total": 100,
	})
	if err != nil {
		t.Fatalf("winPointsDist() failed: %v", err)
	}
	assertScalarProbDist(t, got, scalarProbDist{
		1000: 0.25,
		2000: 0.75,
	})
}

func TestWinPointsDist_ReturnsErrorWithInvalidTotal(t *testing.T) {
	_, err := winPointsDist(map[string]int{
		"1000":  1,
		"total": 0,
	})
	if err == nil {
		t.Fatal("winPointsDist() succeeded unexpectedly")
	}
}

func TestWinPointsDist_ReturnsErrorWithInvalidPointKey(t *testing.T) {
	_, err := winPointsDist(map[string]int{
		"bad":   1,
		"total": 1,
	})
	if err == nil {
		t.Fatal("winPointsDist() succeeded unexpectedly")
	}
}

func TestRandomWinScoreDeltaDist(t *testing.T) {
	got, err := randomWinScoreDeltaDist(1, 0, 0.4, map[string]int{
		"1000":  25,
		"2000":  75,
		"total": 100,
	})
	if err != nil {
		t.Fatalf("randomWinScoreDeltaDist() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{-500, 1000, -250, -250}:  0.10,
		{-1000, 2000, -500, -500}: 0.30,
		{-1000, 1000, 0, 0}:       0.05,
		{-2000, 2000, 0, 0}:       0.15,
		{0, 1000, -1000, 0}:       0.05,
		{0, 2000, -2000, 0}:       0.15,
		{0, 1000, 0, -1000}:       0.05,
		{0, 2000, 0, -2000}:       0.15,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestRandomWinScoreDeltaDistFromStats_SelectsDealerPointFreqs(t *testing.T) {
	got, err := randomWinScoreDeltaDistFromStats(0, 0, stubManueStats{
		numWins:         10,
		numSelfDrawWins: 4,
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
		t.Fatalf("randomWinScoreDeltaDistFromStats() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{2000, -2000.0 / 3.0, -2000.0 / 3.0, -2000.0 / 3.0}: 0.4,
		{2000, -2000, 0, 0}: 0.2,
		{2000, 0, -2000, 0}: 0.2,
		{2000, 0, 0, -2000}: 0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestRandomWinScoreDeltaDistFromStats_SelectsNonDealerPointFreqs(t *testing.T) {
	got, err := randomWinScoreDeltaDistFromStats(1, 0, stubManueStats{
		numWins:         10,
		numSelfDrawWins: 4,
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
		t.Fatalf("randomWinScoreDeltaDistFromStats() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{-500, 1000, -250, -250}: 0.4,
		{-1000, 1000, 0, 0}:      0.2,
		{0, 1000, -1000, 0}:      0.2,
		{0, 1000, 0, -1000}:      0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestRandomWinScoreDeltaDistFromStats_ReturnsErrorWithInvalidNumWins(t *testing.T) {
	_, err := randomWinScoreDeltaDistFromStats(0, 0, stubManueStats{})
	if err == nil {
		t.Fatal("randomWinScoreDeltaDistFromStats() succeeded unexpectedly")
	}
}

func TestWinScoreDeltaDistFromPointsDist(t *testing.T) {
	got, err := winScoreDeltaDistFromPointsDist(1, 0, stubManueStats{
		numWins:         10,
		numSelfDrawWins: 4,
	}, scalarProbDist{
		1000: 0.25,
		2000: 0.75,
	})
	if err != nil {
		t.Fatalf("winScoreDeltaDistFromPointsDist() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{-500, 1000, -250, -250}:  0.10,
		{-1000, 2000, -500, -500}: 0.30,
		{-1000, 1000, 0, 0}:       0.05,
		{-2000, 2000, 0, 0}:       0.15,
		{0, 1000, -1000, 0}:       0.05,
		{0, 2000, -2000, 0}:       0.15,
		{0, 1000, 0, -1000}:       0.05,
		{0, 2000, 0, -2000}:       0.15,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestWinScoreDeltaDistFromPointsDist_ReturnsErrorWithInvalidNumWins(t *testing.T) {
	_, err := winScoreDeltaDistFromPointsDist(0, 0, stubManueStats{}, scalarProbDist{1000: 1})
	if err == nil {
		t.Fatal("winScoreDeltaDistFromPointsDist() succeeded unexpectedly")
	}
}

func TestSelfWinScoreDeltaDistFromEstimate(t *testing.T) {
	got, err := selfWinScoreDeltaDistFromEstimate(1, 0, stubManueStats{
		numWins:         10,
		numSelfDrawWins: 4,
	}, winEstimate{
		pointsDist: scalarProbDist{
			1000: 0.25,
			2000: 0.75,
		},
	})
	if err != nil {
		t.Fatalf("selfWinScoreDeltaDistFromEstimate() failed: %v", err)
	}

	want := scoreDeltaProbDist{
		{-500, 1000, -250, -250}:  0.10,
		{-1000, 2000, -500, -500}: 0.30,
		{-1000, 1000, 0, 0}:       0.05,
		{-2000, 2000, 0, 0}:       0.15,
		{0, 1000, -1000, 0}:       0.05,
		{0, 2000, -2000, 0}:       0.15,
		{0, 1000, 0, -1000}:       0.05,
		{0, 2000, 0, -2000}:       0.15,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestExhaustiveDrawProb(t *testing.T) {
	got, err := exhaustiveDrawProb(stubManueStats{
		turnDistribution:    []float64{0.1, 0.2, 0.3, 0.4},
		exhaustiveDrawRatio: 0.27,
	}, 2.75)
	if err != nil {
		t.Fatalf("exhaustiveDrawProb() failed: %v", err)
	}

	want := 0.27 / 0.7
	if !almostEqual(got, want) {
		t.Errorf("exhaustiveDrawProb() = %v, want %v", got, want)
	}
}

func TestExhaustiveDrawProb_ReturnsErrorWithOutOfRangeTurn(t *testing.T) {
	_, err := exhaustiveDrawProb(stubManueStats{
		turnDistribution:    []float64{0.1},
		exhaustiveDrawRatio: 0.1,
	}, 1)
	if err == nil {
		t.Fatal("exhaustiveDrawProb() succeeded unexpectedly")
	}
}

func TestExhaustiveDrawProbOnSelfNoWin(t *testing.T) {
	got, err := exhaustiveDrawProbOnSelfNoWin(stubManueStats{
		turnDistribution:    []float64{0.25, 0.75},
		exhaustiveDrawRatio: 0.25,
	}, 0)
	if err != nil {
		t.Fatalf("exhaustiveDrawProbOnSelfNoWin() failed: %v", err)
	}

	want := 0.35355339059327373
	if !almostEqual(got, want) {
		t.Errorf("exhaustiveDrawProbOnSelfNoWin() = %v, want %v", got, want)
	}
}

func TestExpectedRemainingTurns(t *testing.T) {
	got, err := expectedRemainingTurns(stubManueStats{
		turnDistribution: []float64{
			0,
			0,
			0,
			0.2,
			0.3,
			0.5,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		},
	}, 3.2)
	if err != nil {
		t.Fatalf("expectedRemainingTurns() failed: %v", err)
	}

	if got != 2 {
		t.Errorf("expectedRemainingTurns() = %v, want 2", got)
	}
}

func TestExpectedRemainingTurns_ReturnsZeroWithoutRemainingTurnProb(t *testing.T) {
	got, err := expectedRemainingTurns(stubManueStats{
		turnDistribution: fullTurnDistribution(0),
	}, 3)
	if err != nil {
		t.Fatalf("expectedRemainingTurns() failed: %v", err)
	}

	if got != 0 {
		t.Errorf("expectedRemainingTurns() = %v, want 0", got)
	}
}

func TestExpectedRemainingTurns_ReturnsErrorWithOutOfRangeTurn(t *testing.T) {
	_, err := expectedRemainingTurns(stubManueStats{
		turnDistribution: fullTurnDistribution(0.1),
	}, 18)
	if err == nil {
		t.Fatal("expectedRemainingTurns() succeeded unexpectedly")
	}
}

func TestTenpaiProb_ReturnsOneWithRiichi(t *testing.T) {
	got, err := tenpaiProb(stubManueStats{}, true, 10, 0)
	if err != nil {
		t.Fatalf("tenpaiProb() failed: %v", err)
	}

	if got != 1.0 {
		t.Errorf("tenpaiProb() = %v, want 1", got)
	}
}

func TestTenpaiProb_ReturnsYamitenRatio(t *testing.T) {
	got, err := tenpaiProb(stubManueStats{
		yamitenCounts: map[string]yamitenCount{
			"10,2": {total: 20, tenpai: 5},
		},
	}, false, 10, 2)
	if err != nil {
		t.Fatalf("tenpaiProb() failed: %v", err)
	}

	want := 0.25
	if got != want {
		t.Errorf("tenpaiProb() = %v, want %v", got, want)
	}
}

func TestTenpaiProb_ReturnsOneWithoutStats(t *testing.T) {
	got, err := tenpaiProb(stubManueStats{}, false, 10, 2)
	if err != nil {
		t.Fatalf("tenpaiProb() failed: %v", err)
	}

	if got != 1.0 {
		t.Errorf("tenpaiProb() = %v, want 1", got)
	}
}

func TestTenpaiProb_ReturnsErrorWithInvalidYamitenCounts(t *testing.T) {
	_, err := tenpaiProb(stubManueStats{
		yamitenCounts: map[string]yamitenCount{
			"10,2": {total: 0, tenpai: 0},
		},
	}, false, 10, 2)
	if err == nil {
		t.Fatal("tenpaiProb() succeeded unexpectedly")
	}
}

func TestDealInExpPts(t *testing.T) {
	got, err := dealInExpPts(stubManueStats{avgWinPts: 5500}, 0.8)
	if err != nil {
		t.Fatalf("dealInExpPts() failed: %v", err)
	}

	want := -1100.0
	if !almostEqual(got, want) {
		t.Errorf("dealInExpPts() = %v, want %v", got, want)
	}
}

func TestDealInExpPts_ReturnsErrorWithInvalidSafeProb(t *testing.T) {
	_, err := dealInExpPts(stubManueStats{avgWinPts: 5500}, 1.1)
	if err == nil {
		t.Fatal("dealInExpPts() succeeded unexpectedly")
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

func TestImmediateDealInScoreDeltaDistFromStats_ReturnsErrorWithInvalidPointFreqs(t *testing.T) {
	_, err := immediateDealInScoreDeltaDistFromStats(1, 0, 2, 0.25, stubManueStats{
		nonDealerWinPointFreqs: map[string]int{
			"1000":  1,
			"total": 0,
		},
	})
	if err == nil {
		t.Fatal("immediateDealInScoreDeltaDistFromStats() succeeded unexpectedly")
	}
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

func TestSafeWinExpPts(t *testing.T) {
	got, err := safeWinExpPts(0.8, 4000)
	if err != nil {
		t.Fatalf("safeWinExpPts() failed: %v", err)
	}

	want := 3200.0
	if got != want {
		t.Errorf("safeWinExpPts() = %v, want %v", got, want)
	}
}

func TestSafeWinExpPts_ReturnsErrorWithInvalidSafeProb(t *testing.T) {
	_, err := safeWinExpPts(-0.1, 4000)
	if err == nil {
		t.Fatal("safeWinExpPts() succeeded unexpectedly")
	}
}

func TestExhaustiveDrawExpPts(t *testing.T) {
	got, err := exhaustiveDrawExpPts(0.8, 0.25, 1500)
	if err != nil {
		t.Fatalf("exhaustiveDrawExpPts() failed: %v", err)
	}

	want := 300.0
	if got != want {
		t.Errorf("exhaustiveDrawExpPts() = %v, want %v", got, want)
	}
}

func TestExhaustiveDrawExpPts_ReturnsErrorWithInvalidProb(t *testing.T) {
	if _, err := exhaustiveDrawExpPts(1.1, 0.25, 1500); err == nil {
		t.Fatal("exhaustiveDrawExpPts() succeeded unexpectedly with invalid safe probability")
	}
	if _, err := exhaustiveDrawExpPts(0.8, -0.1, 1500); err == nil {
		t.Fatal("exhaustiveDrawExpPts() succeeded unexpectedly with invalid exhaustive-draw probability")
	}
}

func TestRemainingRoundEndProbs(t *testing.T) {
	drawProb, othersWinProb, err := remainingRoundEndProbs(0.2, 0.3)
	if err != nil {
		t.Fatalf("remainingRoundEndProbs() failed: %v", err)
	}

	if !almostEqual(drawProb, 0.24) {
		t.Errorf("drawProb = %v, want 0.24", drawProb)
	}
	if !almostEqual(othersWinProb, 0.56) {
		t.Errorf("othersWinProb = %v, want 0.56", othersWinProb)
	}
}

func TestRemainingRoundEndProbs_ReturnsErrorWithInvalidProb(t *testing.T) {
	if _, _, err := remainingRoundEndProbs(-0.1, 0.3); err == nil {
		t.Fatal("remainingRoundEndProbs() succeeded unexpectedly with invalid self win probability")
	}
	if _, _, err := remainingRoundEndProbs(0.2, 1.1); err == nil {
		t.Fatal("remainingRoundEndProbs() succeeded unexpectedly with invalid exhaustive-draw probability")
	}
}

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

func TestNotenExhaustiveDrawTenpaiProb(t *testing.T) {
	got, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"16.25": 30,
			"16.5":  20,
			"16.75": 10,
			"17":    40,
			"17.25": 0,
			"17.5":  0,
		},
	}, 16)
	if err != nil {
		t.Fatalf("notenExhaustiveDrawTenpaiProb() failed: %v", err)
	}

	want := 0.5
	if got != want {
		t.Errorf("notenExhaustiveDrawTenpaiProb() = %v, want %v", got, want)
	}
}

func TestNotenExhaustiveDrawTenpaiProb_UsesFutureTurnsOnly(t *testing.T) {
	got, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"15.75": 1000,
			"16":    1000,
			"16.25": 100,
			"16.5":  0,
			"16.75": 0,
			"17":    0,
			"17.25": 0,
			"17.5":  0,
		},
	}, 16)
	if err != nil {
		t.Fatalf("notenExhaustiveDrawTenpaiProb() failed: %v", err)
	}

	want := 0.5
	if got != want {
		t.Errorf("notenExhaustiveDrawTenpaiProb() = %v, want %v", got, want)
	}
}

func TestNotenExhaustiveDrawTenpaiProb_AllowsExistingZeroFreq(t *testing.T) {
	got, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"16.25": 0,
			"16.5":  100,
			"16.75": 0,
			"17":    0,
			"17.25": 0,
			"17.5":  0,
		},
	}, 16)
	if err != nil {
		t.Fatalf("notenExhaustiveDrawTenpaiProb() failed: %v", err)
	}

	want := 0.5
	if got != want {
		t.Errorf("notenExhaustiveDrawTenpaiProb() = %v, want %v", got, want)
	}
}

func TestNotenExhaustiveDrawTenpaiProb_ReturnsErrorWithMissingTurnFreq(t *testing.T) {
	_, err := notenExhaustiveDrawTenpaiProb(stubManueStats{
		exhaustiveDrawNotenCount: 100,
		exhaustiveDrawTenpaiTurnFreqs: map[string]int{
			"16.25": 100,
		},
	}, 16)
	if err == nil {
		t.Fatal("notenExhaustiveDrawTenpaiProb() succeeded unexpectedly")
	}
}

func TestNotenExhaustiveDrawTenpaiProb_ReturnsErrorWithoutFreqs(t *testing.T) {
	_, err := notenExhaustiveDrawTenpaiProb(stubManueStats{}, 16)
	if err == nil {
		t.Fatal("notenExhaustiveDrawTenpaiProb() succeeded unexpectedly")
	}
}

func TestExhaustiveDrawTenpaiProbs(t *testing.T) {
	got := exhaustiveDrawTenpaiProbs([4]float64{0, 0.25, 0.5, 1}, 0.4)
	want := [4]float64{0.4, 0.55, 0.7, 1}
	if got != want {
		t.Errorf("exhaustiveDrawTenpaiProbs() = %v, want %v", got, want)
	}
}

func TestRyukyokuScoreDelta(t *testing.T) {
	got := ryukyokuScoreDelta([4]bool{true, false, true, false})
	want := scoreDelta{1500, -1500, 1500, -1500}
	if got != want {
		t.Errorf("ryukyokuScoreDelta() = %v, want %v", got, want)
	}
}

func TestExhaustiveDrawScoreDeltaDistFromTenpaiProbs(t *testing.T) {
	got := exhaustiveDrawScoreDeltaDistFromTenpaiProbs([4]float64{1, 0, 0.5, 0})
	want := scoreDeltaProbDist{
		{3000, -1000, -1000, -1000}: 0.5,
		{1500, -1500, 1500, -1500}:  0.5,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestFutureExhaustiveDrawScoreDeltaDist(t *testing.T) {
	got := futureExhaustiveDrawScoreDeltaDist([4]float64{1, 0, 0.5, 0}, 0.5)
	want := exhaustiveDrawScoreDeltaDistFromTenpaiProbs([4]float64{1, 0.5, 0.75, 0.5})
	assertScoreDeltaProbDist(t, got, want)
}

func TestExhaustiveDrawAvgPts(t *testing.T) {
	got := exhaustiveDrawAvgPts(0, [4]float64{1, 0, 0.5, 0})
	want := 2250.0
	if got != want {
		t.Errorf("exhaustiveDrawAvgPts() = %v, want %v", got, want)
	}
}
