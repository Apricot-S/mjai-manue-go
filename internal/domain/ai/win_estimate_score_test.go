package ai

import "testing"

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
	if !almostEqual(got.expectedPoints, 1000) {
		t.Errorf("expectedPoints = %v, want 1000", got.expectedPoints)
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
	if got.expectedPoints != 0 {
		t.Errorf("expectedPoints = %v, want 0", got.expectedPoints)
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
	if !almostEqual(got.expectedPoints, 3000) {
		t.Errorf("expectedPoints = %v, want 3000", got.expectedPoints)
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
	if !almostEqual(got.expectedPoints, 1750) {
		t.Errorf("expectedPoints = %v, want 1750", got.expectedPoints)
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
	if !almostEqual(discardEstimate.expectedPoints, 1000) {
		t.Errorf("discard expectedPoints = %v, want 1000", discardEstimate.expectedPoints)
	}

	riichiEstimate := got["riichi-1m"]
	if riichiEstimate.prob != 1 {
		t.Errorf("riichi prob = %v, want 1", riichiEstimate.prob)
	}
	if riichiEstimate.avgPts != 5000 {
		t.Errorf("riichi avgPts = %v, want 5000", riichiEstimate.avgPts)
	}
	if riichiEstimate.expectedPoints != 5000 {
		t.Errorf("riichi expectedPoints = %v, want 5000", riichiEstimate.expectedPoints)
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

func TestWinEstimateAccumulatorSet_MergeReturnsErrorWithMismatchedKeys(t *testing.T) {
	tests := []struct {
		name  string
		other winEstimateAccumulatorSet
	}{
		{
			name: "unknown key",
			other: winEstimateAccumulatorSet{
				"discard-1m": {numTries: 1},
				"unknown":    {numTries: 1},
			},
		},
		{
			name: "missing key",
			other: winEstimateAccumulatorSet{
				"discard-1m": {numTries: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accumulators := winEstimateAccumulatorSet{
				"discard-1m": {
					numTries:   1,
					totalWins:  1,
					totalPts:   1000,
					pointFreqs: map[float64]int{1000: 1},
				},
				"discard-2m": {
					numTries: 1,
				},
			}

			err := accumulators.merge(tt.other)
			if err == nil {
				t.Fatal("merge() succeeded unexpectedly")
			}
			if accumulators["discard-1m"].numTries != 1 {
				t.Errorf("discard-1m numTries = %v, want unchanged 1", accumulators["discard-1m"].numTries)
			}
			if accumulators["discard-1m"].pointFreqs[1000] != 1 {
				t.Errorf("discard-1m pointFreqs[1000] = %v, want unchanged 1", accumulators["discard-1m"].pointFreqs[1000])
			}
			if accumulators["discard-2m"].numTries != 1 {
				t.Errorf("discard-2m numTries = %v, want unchanged 1", accumulators["discard-2m"].numTries)
			}
		})
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
	if !almostEqual(discard1m.expectedPoints, 7000.0/3.0) {
		t.Errorf("discard-1m expectedPoints = %v, want %v", discard1m.expectedPoints, 7000.0/3.0)
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
	if !almostEqual(discard2m.expectedPoints, 8000.0/3.0) {
		t.Errorf("discard-2m expectedPoints = %v, want %v", discard2m.expectedPoints, 8000.0/3.0)
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
