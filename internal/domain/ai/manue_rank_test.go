package ai

import "testing"

func TestWinProbAgainst(t *testing.T) {
	scoreChanges := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}: 0.25,
		{-1000, 1000, 0, 0}: 0.75,
	})
	winProbs := relativeWinProbTable{
		"3000":  0.9,
		"-1000": 0.4,
	}

	got := winProbAgainst(
		scoreChanges,
		0,
		1,
		25000,
		24000,
		1,
		0,
		winProbs,
	)
	want := 0.25*0.9 + 0.75*0.4
	if got != want {
		t.Errorf("winProbAgainst() = %v, want %v", got, want)
	}
}

func TestWinProbFromRelativeScore_UsesStatsWhenAvailable(t *testing.T) {
	winProbs := relativeWinProbTable{
		"1000": 0.75,
	}

	got := winProbFromRelativeScore(1000, winProbs, 1, 0)
	if got != 0.75 {
		t.Errorf("winProbFromRelativeScore() = %v, want 0.75", got)
	}
}

func TestWinProbFromRelativeScore_FallsBackToStartingDealerOrder(t *testing.T) {
	tests := []struct {
		name          string
		relativeScore float64
		selfPosition  int
		otherPosition int
		want          float64
	}{
		{
			name:          "closer to starting dealer wins tie",
			relativeScore: 0,
			selfPosition:  0,
			otherPosition: 1,
			want:          1,
		},
		{
			name:          "farther from starting dealer loses tie",
			relativeScore: 0,
			selfPosition:  1,
			otherPosition: 0,
			want:          0,
		},
		{
			name:          "positive score wins even when farther",
			relativeScore: 100,
			selfPosition:  1,
			otherPosition: 0,
			want:          1,
		},
		{
			name:          "negative score loses even when closer",
			relativeScore: -100,
			selfPosition:  0,
			otherPosition: 1,
			want:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := winProbFromRelativeScore(tt.relativeScore, nil, tt.selfPosition, tt.otherPosition)
			if got != tt.want {
				t.Errorf("winProbFromRelativeScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
