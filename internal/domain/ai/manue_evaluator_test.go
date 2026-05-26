package ai

import (
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestEvaluateCandidateScore(t *testing.T) {
	scoreChanges := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}: 0.25,
		{-500, 500, 0, 0}:   0.75,
	})
	base := candidateScore{
		winProb:            0.2,
		exhaustiveDrawProb: 0.3,
		shanten:            1,
	}

	got := evaluateCandidateScore(base, scoreChanges, 0, 25000, 0, []rankOpponent{
		{
			id:       1,
			score:    25000,
			position: 1,
			winProbs: relativeWinProbTable{
				"2000":  1.0,
				"-1000": 0.0,
			},
		},
		{
			id:       2,
			score:    24000,
			position: 2,
			winProbs: relativeWinProbTable{
				"2000": 1.0,
				"500":  1.0,
			},
		},
		{
			id:       3,
			score:    26000,
			position: 3,
			winProbs: relativeWinProbTable{
				"0":     0.5,
				"-1500": 0.0,
			},
		},
	})

	if !almostEqual(got.expectedPoints, -125) {
		t.Errorf("expectedPoints = %v, want -125", got.expectedPoints)
	}
	if !almostEqual(got.averageRank, 2.625) {
		t.Errorf("averageRank = %v, want 2.625", got.averageRank)
	}
	if got.winProb != base.winProb {
		t.Errorf("winProb = %v, want %v", got.winProb, base.winProb)
	}
	if got.exhaustiveDrawProb != base.exhaustiveDrawProb {
		t.Errorf("exhaustiveDrawProb = %v, want %v", got.exhaustiveDrawProb, base.exhaustiveDrawProb)
	}
	if got.shanten != base.shanten {
		t.Errorf("shanten = %v, want %v", got.shanten, base.shanten)
	}
}

func TestEvaluateCandidateScoreFromState(t *testing.T) {
	scoreChanges := newScoreDeltaProbDist(map[scoreDelta]float64{
		{1000, -1000, 0, 0}: 0.25,
		{-500, 500, 0, 0}:   0.75,
	})
	base := candidateScore{
		winProb:            0.2,
		exhaustiveDrawProb: 0.3,
		shanten:            1,
	}

	got := evaluateCandidateScoreFromState(
		base,
		scoreChanges,
		stubManueStats{
			relativeWinProbs: map[string]map[string]float64{
				"E1,0,1": {
					"2000":  1.0,
					"-1000": 0.0,
				},
				"E1,0,2": {
					"2000": 1.0,
					"500":  1.0,
				},
				"E1,0,3": {
					"0":     0.5,
					"-1500": 0.0,
				},
			},
		},
		stubRankStateViewer{
			nextRoundWind:  wind.East,
			nextRoundNum:   1,
			scores:         [common.NumPlayers]int{25000, 25000, 24000, 26000},
			startingDealer: seat.MustSeat(0),
		},
		seat.MustSeat(0),
	)

	if !almostEqual(got.expectedPoints, -125) {
		t.Errorf("expectedPoints = %v, want -125", got.expectedPoints)
	}
	if !almostEqual(got.averageRank, 2.625) {
		t.Errorf("averageRank = %v, want 2.625", got.averageRank)
	}
	if got.winProb != base.winProb {
		t.Errorf("winProb = %v, want %v", got.winProb, base.winProb)
	}
	if got.exhaustiveDrawProb != base.exhaustiveDrawProb {
		t.Errorf("exhaustiveDrawProb = %v, want %v", got.exhaustiveDrawProb, base.exhaustiveDrawProb)
	}
	if got.shanten != base.shanten {
		t.Errorf("shanten = %v, want %v", got.shanten, base.shanten)
	}
}

func TestCandidateTotalScoreDeltaDist(t *testing.T) {
	score := candidateScore{
		winProb:            0.2,
		exhaustiveDrawProb: 0.3,
		otherWinProb:       0.4,
	}
	immediateDist := scoreDeltaProbDist{
		{}:            0.75,
		{-1000, 1000}: 0.25,
	}
	selfWinDist := scoreDeltaProbDist{{1000, 0, 0, 0}: 1}
	exhaustiveDrawDist := scoreDeltaProbDist{{0, 1000, 0, 0}: 1}
	otherWinDists := []scoreDeltaProbDist{
		{{0, 0, 1000, 0}: 1},
		{{0, 0, 0, 1000}: 1},
	}

	got := candidateTotalScoreDeltaDist(
		score,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		otherWinDists,
	)

	want := scoreDeltaProbDist{
		{-1000, 1000}:   0.25,
		{1000, 0, 0, 0}: 0.15,
		{0, 1000, 0, 0}: 0.225,
		{0, 0, 1000, 0}: 0.15,
		{0, 0, 0, 1000}: 0.15,
	}
	assertScoreDeltaProbDist(t, got, want)
}

func TestEvaluateCandidateFromComponents(t *testing.T) {
	score := candidateScore{
		shanten: 1,
	}
	dealInEstimates := []dealInEstimate{
		{winnerID: 1, prob: 0.2},
		{winnerID: 2, prob: 0.25},
	}
	immediateDist := scoreDeltaProbDist{
		{}: 1,
	}
	selfWinDist := scoreDeltaProbDist{
		{1000, 0, 0, 0}: 1,
	}
	exhaustiveDrawDist := scoreDeltaProbDist{
		{0, 1000, 0, 0}: 1,
	}
	otherWinDists := []scoreDeltaProbDist{
		{{0, 0, 1000, 0}: 1},
		{{0, 0, 0, 1000}: 1},
	}

	got, err := evaluateCandidateFromComponents(
		score,
		dealInEstimates,
		winEstimate{
			prob:   0.2,
			avgPts: 3900,
		},
		0.25,
		1200,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		otherWinDists,
		stubManueStats{
			relativeWinProbs: map[string]map[string]float64{
				"E1,0,1": {"-1000": 0.5, "0": 0.5, "1000": 0.5},
				"E1,0,2": {"-1000": 0.5, "0": 0.5, "1000": 0.5},
				"E1,0,3": {"-1000": 0.5, "0": 0.5, "1000": 0.5},
			},
		},
		stubRankStateViewer{
			nextRoundWind:  wind.East,
			nextRoundNum:   1,
			scores:         [common.NumPlayers]int{25000, 25000, 25000, 25000},
			startingDealer: seat.MustSeat(0),
		},
		seat.MustSeat(0),
	)
	if err != nil {
		t.Fatalf("evaluateCandidateFromComponents() failed: %v", err)
	}

	if !almostEqual(got.dealInProb, 0.4) {
		t.Errorf("dealInProb = %v, want 0.4", got.dealInProb)
	}
	if !almostEqual(got.exhaustiveDrawProb, 0.2) {
		t.Errorf("exhaustiveDrawProb = %v, want 0.2", got.exhaustiveDrawProb)
	}
	if !almostEqual(got.otherWinProb, 0.6) {
		t.Errorf("otherWinProb = %v, want 0.6", got.otherWinProb)
	}
	if got.averageWinPoints != 3900 {
		t.Errorf("averageWinPoints = %v, want 3900", got.averageWinPoints)
	}
	if got.exhaustiveDrawAveragePoints != 1200 {
		t.Errorf("exhaustiveDrawAveragePoints = %v, want 1200", got.exhaustiveDrawAveragePoints)
	}
	if !almostEqual(got.expectedPoints, 200) {
		t.Errorf("expectedPoints = %v, want 200", got.expectedPoints)
	}
	if !almostEqual(got.averageRank, 2.5) {
		t.Errorf("averageRank = %v, want 2.5", got.averageRank)
	}
}

func TestEvaluateCandidateFromComponents_ReturnsErrorWithInvalidEstimate(t *testing.T) {
	_, err := evaluateCandidateFromComponents(
		candidateScore{},
		[]dealInEstimate{{winnerID: 1, prob: 1.1}},
		winEstimate{},
		0.25,
		0,
		nil,
		nil,
		nil,
		nil,
		stubManueStats{},
		stubRankStateViewer{},
		seat.MustSeat(0),
	)
	if err == nil {
		t.Fatal("evaluateCandidateFromComponents() succeeded unexpectedly")
	}
}

func TestManueAgent_evaluateActionCandidate_ReturnsErrorWithMissingWinEstimate(t *testing.T) {
	context := actionEvaluationContext{
		stats:        validStubManueStats(),
		state:        stubCandidateEvaluationStateViewer{},
		self:         seat.MustSeat(0),
		winEstimates: map[string]winEstimate{},
	}

	_, err := newTestManueAgent(t, 0).evaluateActionCandidate(context, actionCandidate{
		traceKey:    "-1.1m",
		discardTile: tile.MustTileFromCode("1m"),
	})
	if err == nil {
		t.Fatal("evaluateActionCandidate() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "missing win estimate") {
		t.Errorf("evaluateActionCandidate() error = %v, want missing win estimate", err)
	}
}
