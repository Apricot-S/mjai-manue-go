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
