package ai

import (
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestEvaluateCandidateFromComponents(t *testing.T) {
	dealInEstimates := []dealInEstimate{
		{winnerID: 1, prob: 0.2},
		{winnerID: 2, prob: 0.25},
	}
	immediateDist := scoreDeltaProbDist{
		{}:            0.75,
		{-1000, 1000}: 0.25,
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
		dealInEstimates,
		winEstimate{
			prob:   0.2,
			avgPts: 3900,
		},
		0.375,
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
	if !almostEqual(got.exhaustiveDrawProb, 0.3) {
		t.Errorf("exhaustiveDrawProb = %v, want 0.3", got.exhaustiveDrawProb)
	}
	if !almostEqual(got.otherWinProb, 0.5) {
		t.Errorf("otherWinProb = %v, want 0.5", got.otherWinProb)
	}
	if got.averageWinPoints != 3900 {
		t.Errorf("averageWinPoints = %v, want 3900", got.averageWinPoints)
	}
	if got.exhaustiveDrawAveragePoints != 1200 {
		t.Errorf("exhaustiveDrawAveragePoints = %v, want 1200", got.exhaustiveDrawAveragePoints)
	}
	if !almostEqual(got.expectedPoints, -100) {
		t.Errorf("expectedPoints = %v, want -100", got.expectedPoints)
	}
	if !almostEqual(got.averageRank, 2.625) {
		t.Errorf("averageRank = %v, want 2.625", got.averageRank)
	}
}

func TestEvaluateCandidateFromComponents_ReturnsErrorWithInvalidEstimate(t *testing.T) {
	_, err := evaluateCandidateFromComponents(
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

func TestCandidateEvaluator_evaluateCandidate_ReturnsErrorWithMissingWinEstimate(t *testing.T) {
	context := candidateEvaluationContext{
		stats:        validStubManueStats(),
		state:        stubCandidateEvaluationStateViewer{},
		self:         seat.MustSeat(0),
		winEstimates: map[string]winEstimate{},
	}

	evaluator := candidateEvaluator{
		stats:  validStubManueStats(),
		danger: stubDangerEstimator{},
	}
	_, err := evaluator.evaluateCandidate(context, actionCandidate{
		traceKey:    "-1.1m",
		discardTile: tile.MustTileFromCode("1m"),
	})
	if err == nil {
		t.Fatal("evaluateCandidate() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "missing win estimate") {
		t.Errorf("evaluateCandidate() error = %v, want missing win estimate", err)
	}
}

func TestCandidateEvaluator_newEvaluationContext_ReturnsErrorWithInvalidTrialCount(t *testing.T) {
	discard := tile.MustTileFromCode("1m")
	var throwable hand.TileCounts34
	throwable[discard.ID()] = 1
	evaluator := candidateEvaluator{
		stats:  validStubManueStats(),
		danger: stubDangerEstimator{},
		rng:    rand.New(rand.NewPCG(0, 0)),
		trials: 0,
	}

	_, err := evaluator.newEvaluationContext(
		stubCandidateEvaluationStateViewer{
			roundWind:    wind.East,
			seatWinds:    [common.NumPlayers]wind.Wind{wind.East, wind.South, wind.West, wind.North},
			dealer:       seat.MustSeat(0),
			numLeftTiles: common.NumPlayers * 4,
			players: [common.NumPlayers]player.PlayerViewer{
				stubPlayerViewer{
					hand: visibleHandFromCodes("1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "2p", "3p", "4p"),
				},
				stubPlayerViewer{},
				stubPlayerViewer{},
				stubPlayerViewer{},
			},
		},
		seat.MustSeat(0),
		[]actionCandidate{{
			traceKey:    "-1.1m",
			discardTile: discard,
			turnHand:    visibleHandFromCodes("1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1p", "2p", "3p", "4p"),
			shantenGoals: []service.Goal{{
				Shanten:         0,
				ThrowableVector: throwable,
			}},
		}},
	)
	if err == nil {
		t.Fatal("newEvaluationContext() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "numTries must be positive") {
		t.Errorf("newEvaluationContext() error = %v, want numTries validation", err)
	}
}

func TestCandidateEvaluator_dealInEstimates_SafeTileHasZeroDealInProb(t *testing.T) {
	self := seat.MustSeat(0)
	discard := tile.MustTileFromCode("5m")
	state := stubCandidateEvaluationStateViewer{
		roundWind:    wind.East,
		seatWinds:    [common.NumPlayers]wind.Wind{wind.East, wind.South, wind.West, wind.North},
		dealer:       self,
		numLeftTiles: 4,
		safeTiles: [common.NumPlayers]tile.Tiles{
			nil,
			{discard},
			{discard},
			{discard},
		},
		players: [common.NumPlayers]player.PlayerViewer{
			stubPlayerViewer{},
			stubPlayerViewer{},
			stubPlayerViewer{},
			stubPlayerViewer{},
		},
	}
	evaluator := candidateEvaluator{
		stats:  validStubManueStats(),
		danger: NewDangerEstimator(stubDangerTreeLeaf{prob: 0.75}),
	}

	got, err := evaluator.dealInEstimates(state, self, discard)
	if err != nil {
		t.Fatalf("dealInEstimates() failed: %v", err)
	}
	if len(got) != common.NumPlayers-1 {
		t.Fatalf("len(dealInEstimates()) = %d, want %d", len(got), common.NumPlayers-1)
	}
	for _, estimate := range got {
		if estimate.prob != 0 {
			t.Errorf("deal-in prob for winner %d = %v, want 0 for safe tile", estimate.winnerID, estimate.prob)
		}
	}
}

type stubDangerTreeLeaf struct {
	prob float64
}

func (s stubDangerTreeLeaf) LeafProb() (float64, bool) {
	return s.prob, true
}

func (s stubDangerTreeLeaf) Feature() (string, bool) {
	return "", false
}

func (s stubDangerTreeLeaf) NegativeNode() DangerTreeNode {
	return nil
}

func (s stubDangerTreeLeaf) PositiveNode() DangerTreeNode {
	return nil
}

func visibleHandFromCodes(codes ...string) *hand.VisibleHand {
	tiles := make([]tile.Tile, 0, len(codes))
	for _, code := range codes {
		tiles = append(tiles, tile.MustTileFromCode(code))
	}
	return hand.MustVisibleHand(tiles)
}
