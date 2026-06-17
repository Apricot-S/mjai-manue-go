package ai

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type stubRankStateViewer struct {
	nextRoundWind  wind.Wind
	nextRoundNum   int
	scores         [common.NumPlayers]int
	startingDealer seat.Seat
}

func (s stubRankStateViewer) NextRound() (wind.Wind, int) {
	return s.nextRoundWind, s.nextRoundNum
}

func (s stubRankStateViewer) Scores() [common.NumPlayers]int {
	return s.scores
}

func (s stubRankStateViewer) StartingDealer() seat.Seat {
	return s.startingDealer
}

func TestAverageRank(t *testing.T) {
	scoreChanges := newScoreDeltaProbDist(map[scoreDelta]float64{
		{}: 1,
	})
	opponents := []rankOpponent{
		{
			id:       1,
			score:    24000,
			position: 1,
			winProbs: relativeWinProbTable{
				"1000": 0.8,
			},
		},
		{
			id:       2,
			score:    26000,
			position: 2,
			winProbs: relativeWinProbTable{
				"-1000": 0.4,
			},
		},
		{
			id:       3,
			score:    25000,
			position: 3,
			winProbs: relativeWinProbTable{
				"0": 0.6,
			},
		},
	}

	got := averageRank(scoreChanges, 0, 25000, 0, opponents)
	want := 4.0 - (0.8 + 0.4 + 0.6)
	if !almostEqual(got, want) {
		t.Errorf("averageRank() = %v, want %v", got, want)
	}
}

func TestBuildRankOpponents(t *testing.T) {
	got := buildRankOpponents(stubManueStats{
		relativeWinProbs: map[string]map[string]float64{
			"S2,0,1": {"1000": 0.6},
			"S2,0,2": {"1000": 0.7},
			"S2,0,3": {"1000": 0.8},
		},
	}, stubRankStateViewer{
		nextRoundWind:  wind.South,
		nextRoundNum:   2,
		scores:         [common.NumPlayers]int{27000, 24000, 26000, 23000},
		startingDealer: seat.MustSeat(1),
	}, seat.MustSeat(1))

	if len(got) != 3 {
		t.Fatalf("len(buildRankOpponents()) = %d, want 3", len(got))
	}
	tests := []struct {
		index        int
		wantID       int
		wantScore    float64
		wantPosition int
		wantWinProb  float64
	}{
		{index: 0, wantID: 0, wantScore: 27000, wantPosition: 3, wantWinProb: 0.8},
		{index: 1, wantID: 2, wantScore: 26000, wantPosition: 1, wantWinProb: 0.6},
		{index: 2, wantID: 3, wantScore: 23000, wantPosition: 2, wantWinProb: 0.7},
	}
	for _, tt := range tests {
		opponent := got[tt.index]
		if opponent.id != tt.wantID {
			t.Errorf("opponents[%d].id = %d, want %d", tt.index, opponent.id, tt.wantID)
		}
		if opponent.score != tt.wantScore {
			t.Errorf("opponents[%d].score = %v, want %v", tt.index, opponent.score, tt.wantScore)
		}
		if opponent.position != tt.wantPosition {
			t.Errorf("opponents[%d].position = %d, want %d", tt.index, opponent.position, tt.wantPosition)
		}
		if opponent.winProbs["1000"] != tt.wantWinProb {
			t.Errorf("opponents[%d].winProbs[\"1000\"] = %v, want %v", tt.index, opponent.winProbs["1000"], tt.wantWinProb)
		}
	}
}

func TestRelativeWinProbs(t *testing.T) {
	got := relativeWinProbs(stubManueStats{
		relativeWinProbs: map[string]map[string]float64{
			"E1,0,1": {
				"0": 0.5,
			},
		},
	}, wind.East, 1, 0, 1)

	if got["0"] != 0.5 {
		t.Errorf("relativeWinProbs()[\"0\"] = %v, want 0.5", got["0"])
	}
}

func TestRelativeWinProbs_ReturnsNilWithoutStats(t *testing.T) {
	got := relativeWinProbs(stubManueStats{}, wind.East, 1, 0, 1)
	if got != nil {
		t.Errorf("relativeWinProbs() = %v, want nil", got)
	}
}

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
