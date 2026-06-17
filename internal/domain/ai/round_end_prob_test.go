package ai

import (
	"math"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
)

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

	want := math.Pow(0.25, float64(common.NumPlayers-1)/float64(common.NumPlayers))
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

func TestExpectedRemainingTurns_ReturnsZeroAtFinalTurn(t *testing.T) {
	got, err := expectedRemainingTurns(stubManueStats{
		turnDistribution: fullTurnDistribution(0.1),
	}, 17.5)
	if err != nil {
		t.Fatalf("expectedRemainingTurns() failed: %v", err)
	}

	if got != 0 {
		t.Errorf("expectedRemainingTurns() = %v, want 0", got)
	}
}

func TestExpectedRemainingTurns_ReturnsErrorWithOutOfRangeTurn(t *testing.T) {
	tests := []struct {
		name        string
		currentTurn float64
	}{
		{
			name:        "negative",
			currentTurn: -0.25,
		},
		{
			name:        "after final turn",
			currentTurn: 17.75,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := expectedRemainingTurns(stubManueStats{
				turnDistribution: fullTurnDistribution(0.1),
			}, tt.currentTurn)
			if err == nil {
				t.Fatal("expectedRemainingTurns() succeeded unexpectedly")
			}
		})
	}
}

func TestRemainingRoundEndProbs(t *testing.T) {
	exhaustiveDrawProb, otherWinProb, err := remainingRoundEndProbs(0.2, 0.3)
	if err != nil {
		t.Fatalf("remainingRoundEndProbs() failed: %v", err)
	}

	if !almostEqual(exhaustiveDrawProb, 0.24) {
		t.Errorf("exhaustiveDrawProb = %v, want 0.24", exhaustiveDrawProb)
	}
	if !almostEqual(otherWinProb, 0.56) {
		t.Errorf("otherWinProb = %v, want 0.56", otherWinProb)
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
