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

func TestRandomWinScoreDeltaDist_SelectsDealerPointFreqs(t *testing.T) {
	got := randomWinScoreDeltaDist(0, 0, stubManueStats{
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

	want := scoreDeltaProbDist{
		{2000, -2000.0 / 3.0, -2000.0 / 3.0, -2000.0 / 3.0}: 0.4,
		{2000, -2000, 0, 0}: 0.2,
		{2000, 0, -2000, 0}: 0.2,
		{2000, 0, 0, -2000}: 0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestRandomWinScoreDeltaDist_SelectsNonDealerPointFreqs(t *testing.T) {
	got := randomWinScoreDeltaDist(1, 0, stubManueStats{
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

	want := scoreDeltaProbDist{
		{-500, 1000, -250, -250}: 0.4,
		{-1000, 1000, 0, 0}:      0.2,
		{0, 1000, -1000, 0}:      0.2,
		{0, 1000, 0, -1000}:      0.2,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestWinScoreDeltaDist(t *testing.T) {
	got := winScoreDeltaDist(1, 0, stubManueStats{
		numWins:         10,
		numSelfDrawWins: 4,
	}, scalarProbDist{
		1000: 0.25,
		2000: 0.75,
	})

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
