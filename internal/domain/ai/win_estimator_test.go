package ai

import (
	"math/rand/v2"
	"reflect"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestCandidateTraceKeys(t *testing.T) {
	got, err := candidateTraceKeys([]actionCandidate{
		{traceKey: "-1.5m"},
		{traceKey: "0.5m"},
	})
	if err != nil {
		t.Fatalf("candidateTraceKeys() failed: %v", err)
	}

	want := []string{"-1.5m", "0.5m"}
	if len(got) != len(want) {
		t.Fatalf("len(candidateTraceKeys()) = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("candidateTraceKeys()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCandidateTraceKeys_ReturnsErrorWithInvalidKey(t *testing.T) {
	tests := []struct {
		name       string
		candidates []actionCandidate
	}{
		{
			name:       "empty",
			candidates: []actionCandidate{{traceKey: ""}},
		},
		{
			name: "duplicate",
			candidates: []actionCandidate{
				{traceKey: "-1.5m"},
				{traceKey: "-1.5m"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := candidateTraceKeys(tt.candidates)
			if err == nil {
				t.Fatal("candidateTraceKeys() succeeded unexpectedly")
			}
		})
	}
}

func TestFilteredWinEstimateGoals(t *testing.T) {
	tests := []struct {
		name      string
		candidate actionCandidate
		want      []int
	}{
		{
			name: "keeps all goals normally",
			candidate: actionCandidate{
				discardTile: tile.MustTileFromCode("1m"),
				shanten:     2,
				shantenGoals: []service.Goal{
					{Shanten: 1, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 2, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 3, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "riichi keeps only ready goals",
			candidate: actionCandidate{
				riichi:        true,
				scoreAsRiichi: true,
				discardTile:   tile.MustTileFromCode("1m"),
				shanten:       0,
				shantenGoals: []service.Goal{
					{Shanten: 0, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 1, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{0},
		},
		{
			name: "heavy shanten drops extra tile goals",
			candidate: actionCandidate{
				discardTile: tile.MustTileFromCode("1m"),
				baseShanten: 4,
				shanten:     1,
				shantenGoals: []service.Goal{
					{Shanten: 3, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 4, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 5, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{3, 4},
		},
		{
			name: "riichi and heavy shanten combine",
			candidate: actionCandidate{
				riichi:        true,
				scoreAsRiichi: true,
				discardTile:   tile.MustTileFromCode("1m"),
				baseShanten:   4,
				shanten:       1,
				shantenGoals: []service.Goal{
					{Shanten: 0, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 4, ThrowableVector: hand.TileCounts34{0: 1}},
					{Shanten: 5, ThrowableVector: hand.TileCounts34{0: 1}},
				},
			},
			want: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filteredWinEstimateGoals(tt.candidate)
			if len(got) != len(tt.want) {
				t.Fatalf("len(filteredWinEstimateGoals()) = %d, want %d", len(got), len(tt.want))
			}
			for i, want := range tt.want {
				if got[i].Shanten != want {
					t.Errorf("filteredWinEstimateGoals()[%d].Shanten = %d, want %d", i, got[i].Shanten, want)
				}
			}
		})
	}
}

func TestScoredWinEstimateGoals(t *testing.T) {
	turnHand := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
		tile.MustTileFromCode("2s"),
	})
	afterDiscardHand, err := turnHand.Discard(tile.MustTileFromCode("1m"))
	if err != nil {
		t.Fatalf("Discard(1m) failed: %v", err)
	}
	candidate := actionCandidate{
		discardTile:      tile.MustTileFromCode("1m"),
		turnHand:         turnHand,
		afterDiscardHand: afterDiscardHand,
		shantenGoals: []service.Goal{
			{
				Blocks: []block.Block{
					block.MustSequence(tile.MustTileFromCode("2m")),
					block.MustSequence(tile.MustTileFromCode("3p")),
					block.MustSequence(tile.MustTileFromCode("4s")),
					block.MustSequence(tile.MustTileFromCode("6s")),
					block.MustPair(tile.MustTileFromCode("2s")),
				},
				RequiredVector:  hand.TileCounts34{19: 1},
				ThrowableVector: hand.TileCounts34{0: 1},
			},
			{
				Blocks: []block.Block{
					block.MustSequence(tile.MustTileFromCode("1m")),
					block.MustSequence(tile.MustTileFromCode("2p")),
					block.MustSequence(tile.MustTileFromCode("3s")),
					block.MustSequence(tile.MustTileFromCode("6s")),
					block.MustPair(tile.MustTileFromCode("E")),
				},
				RequiredVector:  hand.TileCounts34{1: 1},
				ThrowableVector: hand.TileCounts34{0: 1},
			},
		},
	}

	got, err := scoredWinEstimateGoals(candidate, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
		dealer:    false,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoals() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(scoredWinEstimateGoals()) = %d, want 1", len(got))
	}
	wantPoints := service.RonPoints(30, 2, false)
	if got[0].points != float64(wantPoints) {
		t.Errorf("scoredWinEstimateGoals()[0].points = %v, want %v", got[0].points, wantPoints)
	}
	if got[0].RequiredVector[19] != 1 {
		t.Errorf("scoredWinEstimateGoals()[0].RequiredVector[19] = %d, want 1", got[0].RequiredVector[19])
	}
}

func TestScoredWinEstimateGoalsBuildsHandFromGoalBlocks(t *testing.T) {
	candidate := actionCandidate{
		discardTile: tile.MustTileFromCode("?"),
		turnHand: hand.MustVisibleHand([]tile.Tile{
			tile.MustTileFromCode("5mr"),
			tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("2p"),
			tile.MustTileFromCode("3p"),
			tile.MustTileFromCode("4p"),
			tile.MustTileFromCode("3s"),
			tile.MustTileFromCode("4s"),
			tile.MustTileFromCode("6s"),
			tile.MustTileFromCode("6s"),
			tile.MustTileFromCode("6s"),
			tile.MustTileFromCode("9s"),
			tile.MustTileFromCode("9s"),
		}),
		scoreAsRiichi: true,
		shantenGoals: []service.Goal{
			{
				Blocks: []block.Block{
					block.MustTriplet(tile.MustTileFromCode("1m")),
					block.MustSequence(tile.MustTileFromCode("2p")),
					block.MustSequence(tile.MustTileFromCode("3s")),
					block.MustTriplet(tile.MustTileFromCode("6s")),
					block.MustPair(tile.MustTileFromCode("9s")),
				},
			},
		},
	}

	got, err := scoredWinEstimateGoals(candidate, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
		dealer:    false,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoals() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(scoredWinEstimateGoals()) = %d, want 1", len(got))
	}
	wantPoints := service.RonPoints(40, 1, false)
	if got[0].points != float64(wantPoints) {
		t.Errorf("scoredWinEstimateGoals()[0].points = %v, want %v", got[0].points, wantPoints)
	}
}

func TestScoredWinEstimateGoalsRequiresTurnHand(t *testing.T) {
	_, err := scoredWinEstimateGoals(actionCandidate{}, winEstimateGoalContext{})
	if err == nil {
		t.Fatal("scoredWinEstimateGoals() succeeded unexpectedly")
	}
}

func TestScoredWinEstimateGoalsByKey(t *testing.T) {
	turnHand := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4m"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4p"),
		tile.MustTileFromCode("5p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("5s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("6s"),
		tile.MustTileFromCode("7s"),
		tile.MustTileFromCode("8s"),
		tile.MustTileFromCode("2s"),
	})
	goals := []service.Goal{
		{
			Blocks: []block.Block{
				block.MustSequence(tile.MustTileFromCode("2m")),
				block.MustSequence(tile.MustTileFromCode("3p")),
				block.MustSequence(tile.MustTileFromCode("4s")),
				block.MustSequence(tile.MustTileFromCode("6s")),
				block.MustPair(tile.MustTileFromCode("2s")),
			},
			RequiredVector:  hand.TileCounts34{19: 1},
			ThrowableVector: hand.TileCounts34{0: 1},
		},
	}
	afterDiscard1m, err := turnHand.Discard(tile.MustTileFromCode("1m"))
	if err != nil {
		t.Fatalf("Discard(1m) failed: %v", err)
	}
	afterDiscard2m, err := turnHand.Discard(tile.MustTileFromCode("2m"))
	if err != nil {
		t.Fatalf("Discard(2m) failed: %v", err)
	}
	candidates := []actionCandidate{
		{
			traceKey:         "-1.1m",
			discardTile:      tile.MustTileFromCode("1m"),
			turnHand:         turnHand,
			afterDiscardHand: afterDiscard1m,
			shantenGoals:     goals,
		},
		{
			traceKey:         "-1.2m",
			discardTile:      tile.MustTileFromCode("2m"),
			turnHand:         turnHand,
			afterDiscardHand: afterDiscard2m,
		},
	}

	got, err := scoredWinEstimateGoalsByKey(candidates, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoalsByKey() failed: %v", err)
	}
	if len(got["-1.1m"]) != 1 {
		t.Fatalf("len(scored goals for -1.1m) = %d, want 1", len(got["-1.1m"]))
	}
	if got["-1.1m"][0].points != float64(service.RonPoints(30, 2, false)) {
		t.Errorf("points for -1.1m = %v, want %v", got["-1.1m"][0].points, service.RonPoints(30, 2, false))
	}
	if goals, ok := got["-1.2m"]; !ok {
		t.Fatal("scored goals for -1.2m missing")
	} else if len(goals) != 0 {
		t.Errorf("len(scored goals for -1.2m) = %d, want 0", len(goals))
	}
}

func TestScoredWinEstimateGoalsByKeyUsesCandidateMelds(t *testing.T) {
	turnHand := hand.MustVisibleHand([]tile.Tile{
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("3p"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("4s"),
		tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("9s"),
	})
	goals := []service.Goal{
		{
			Blocks: []block.Block{
				block.MustTriplet(tile.MustTileFromCode("2m")),
				block.MustTriplet(tile.MustTileFromCode("3p")),
				block.MustTriplet(tile.MustTileFromCode("4s")),
				block.MustPair(tile.MustTileFromCode("9s")),
			},
			RequiredVector: hand.TileCounts34{0: 1},
		},
	}
	dragonPon := meld.MustPon(
		tile.MustTileFromCode("F"),
		[2]tile.Tile{tile.MustTileFromCode("F"), tile.MustTileFromCode("F")},
		seat.MustSeat(1),
	)
	candidates := []actionCandidate{
		{
			traceKey:         "none",
			discardTile:      tile.MustTileFromCode("?"),
			turnHand:         turnHand,
			afterDiscardHand: turnHand,
			shantenGoals:     goals,
		},
		{
			traceKey:         "0.1m",
			discardTile:      tile.MustTileFromCode("?"),
			turnHand:         turnHand,
			afterDiscardHand: turnHand,
			melds:            []meld.Meld{dragonPon},
			shantenGoals:     goals,
		},
	}

	got, err := scoredWinEstimateGoalsByKey(candidates, winEstimateGoalContext{
		roundWind: wind.East,
		seatWind:  wind.South,
	})
	if err != nil {
		t.Fatalf("scoredWinEstimateGoalsByKey() failed: %v", err)
	}

	if len(got["none"]) != 1 {
		t.Fatalf("len(scored goals for none) = %d, want 1", len(got["none"]))
	}
	if len(got["0.1m"]) != 1 {
		t.Fatalf("len(scored goals for 0.1m) = %d, want 1", len(got["0.1m"]))
	}
	want := float64(service.RonPoints(30, 3, false))
	if got["0.1m"][0].points != want {
		t.Errorf("points for 0.1m = %v, want %v", got["0.1m"][0].points, want)
	}
	if got["0.1m"][0].points <= got["none"][0].points {
		t.Errorf("points with candidate melds = %v, want greater than %v", got["0.1m"][0].points, got["none"][0].points)
	}
}

func TestScoredWinEstimateGoalsByKeyWrapsCandidateError(t *testing.T) {
	_, err := scoredWinEstimateGoalsByKey(
		[]actionCandidate{{traceKey: "-1.1m"}},
		winEstimateGoalContext{},
	)
	if err == nil {
		t.Fatal("scoredWinEstimateGoalsByKey() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "-1.1m") {
		t.Errorf("scoredWinEstimateGoalsByKey() error = %v, want trace key", err)
	}
}

func TestTrialTileCounts(t *testing.T) {
	got := trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("E"),
	})

	if got[4] != 2 {
		t.Errorf("5m count = %d, want 2", got[4])
	}
	if got[27] != 1 {
		t.Errorf("E count = %d, want 1", got[27])
	}
}

func TestWallTilesFromCounts(t *testing.T) {
	got, err := wallTilesFromCounts(hand.TileCounts34{
		0:  2,
		4:  1,
		27: 1,
	})
	if err != nil {
		t.Fatalf("wallTilesFromCounts() failed: %v", err)
	}
	want := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("E"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("wallTilesFromCounts() = %v, want %v", got, want)
	}
}

func TestWallTilesFromCountsRejectsNegativeCount(t *testing.T) {
	_, err := wallTilesFromCounts(hand.TileCounts34{0: -1})
	if err == nil {
		t.Fatal("wallTilesFromCounts() succeeded unexpectedly")
	}
}

func TestUnseenWallFromVisibleTiles(t *testing.T) {
	got, err := unseenWallFromVisibleTiles([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("E"),
	})
	if err != nil {
		t.Fatalf("unseenWallFromVisibleTiles() failed: %v", err)
	}
	counts := trialTileCounts(got)
	if counts[0] != 3 {
		t.Errorf("1m unseen count = %d, want 3", counts[0])
	}
	if counts[4] != 3 {
		t.Errorf("5m unseen count = %d, want 3", counts[4])
	}
	if counts[27] != 3 {
		t.Errorf("E unseen count = %d, want 3", counts[27])
	}
	if numTiles := (&counts).NumTiles(); numTiles != 133 {
		t.Errorf("unseen wall tile count = %d, want 133", numTiles)
	}
}

func TestUnseenWallFromVisibleTilesRejectsInvalidVisibleTiles(t *testing.T) {
	if _, err := unseenWallFromVisibleTiles([]tile.Tile{tile.MustTileFromCode("?")}); err == nil {
		t.Fatal("unseenWallFromVisibleTiles() accepted unknown tile")
	}

	_, err := unseenWallFromVisibleTiles([]tile.Tile{
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
	})
	if err == nil {
		t.Fatal("unseenWallFromVisibleTiles() accepted tile visible more than 4 times")
	}
}

func TestTrialTilesFromWall(t *testing.T) {
	wall := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
		tile.MustTileFromCode("3m"),
	}

	got, err := trialTilesFromWall(wall, 2)
	if err != nil {
		t.Fatalf("trialTilesFromWall() failed: %v", err)
	}
	want := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("2m"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("trialTilesFromWall() = %v, want %v", got, want)
	}

	got[0] = tile.MustTileFromCode("9m")
	if wall[0] != tile.MustTileFromCode("1m") {
		t.Errorf("wall[0] = %s, want unchanged 1m", wall[0])
	}
}

func TestTrialTilesFromWallRejectsInvalidNumDraws(t *testing.T) {
	wall := []tile.Tile{tile.MustTileFromCode("1m")}
	if _, err := trialTilesFromWall(wall, -1); err == nil {
		t.Fatal("trialTilesFromWall(-1) succeeded unexpectedly")
	}
	if _, err := trialTilesFromWall(wall, 2); err == nil {
		t.Fatal("trialTilesFromWall(2) succeeded unexpectedly")
	}
}

func TestCanAchieveGoalWithTrialTiles(t *testing.T) {
	goal := service.Goal{
		RequiredVector: hand.TileCounts34{
			0:  1,
			4:  2,
			27: 1,
		},
	}

	if !canAchieveGoalWithTrialTiles(goal, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("5mr"),
		tile.MustTileFromCode("E"),
	})) {
		t.Fatal("canAchieveGoalWithTrialTiles() = false, want true")
	}
	if canAchieveGoalWithTrialTiles(goal, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
		tile.MustTileFromCode("E"),
	})) {
		t.Fatal("canAchieveGoalWithTrialTiles() = true, want false")
	}
}

func TestTrialWinPts(t *testing.T) {
	goals := []winEstimateGoal{
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{0: 1},
			},
			points: 1000,
		},
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{0: 1, 4: 1},
			},
			points: 2000,
		},
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{27: 1},
			},
			points: 8000,
		},
	}

	got, ok, err := trialWinPts(goals, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}))
	if err != nil {
		t.Fatalf("trialWinPts() failed: %v", err)
	}
	if !ok {
		t.Fatal("trialWinPts() ok = false, want true")
	}
	if got != 2000 {
		t.Errorf("trialWinPts() = %v, want 2000", got)
	}

	got, ok, err = trialWinPts(goals, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("2m"),
	}))
	if err != nil {
		t.Fatalf("trialWinPts() failed: %v", err)
	}
	if ok {
		t.Fatal("trialWinPts() ok = true, want false")
	}
	if got != 0 {
		t.Errorf("trialWinPts() = %v, want 0", got)
	}
}

func TestTrialWinPtsRejectsNonPositivePoints(t *testing.T) {
	_, _, err := trialWinPts([]winEstimateGoal{
		{
			Goal: service.Goal{
				RequiredVector: hand.TileCounts34{0: 1},
			},
			points: 0,
		},
	}, trialTileCounts([]tile.Tile{tile.MustTileFromCode("1m")}))
	if err == nil {
		t.Fatal("trialWinPts() succeeded unexpectedly")
	}
}

func TestCandidateTrialWinPts(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
		{traceKey: "-1.2m"},
		{traceKey: "0.3m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1},
				},
				points: 1000,
			},
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1, 4: 1},
				},
				points: 3900,
			},
		},
		"-1.2m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{27: 1},
				},
				points: 8000,
			},
		},
		"0.3m": {},
	}

	got, err := candidateTrialWinPts(candidates, goalsByKey, trialTileCounts([]tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}))
	if err != nil {
		t.Fatalf("candidateTrialWinPts() failed: %v", err)
	}
	want := map[string]float64{
		"-1.1m": 3900,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("candidateTrialWinPts() = %#v, want %#v", got, want)
	}
}

func TestCandidateTrialWinPtsRequiresGoalsForEveryCandidate(t *testing.T) {
	_, err := candidateTrialWinPts(
		[]actionCandidate{{traceKey: "-1.1m"}},
		map[string][]winEstimateGoal{},
		hand.TileCounts34{},
	)
	if err == nil {
		t.Fatal("candidateTrialWinPts() succeeded unexpectedly")
	}
}

func TestWinEstimatesFromShuffledWall(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1, 4: 1},
				},
				points: 3900,
			},
		},
	}
	wall := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}

	got, err := winEstimatesFromShuffledWall(
		candidates,
		goalsByKey,
		wall,
		2,
		3,
		rand.New(rand.NewPCG(1, 0)),
	)
	if err != nil {
		t.Fatalf("winEstimatesFromShuffledWall() failed: %v", err)
	}

	estimate := got["-1.1m"]
	if estimate.prob != 1 {
		t.Errorf("estimate.prob = %v, want 1", estimate.prob)
	}
	if estimate.avgPts != 3900 {
		t.Errorf("estimate.avgPts = %v, want 3900", estimate.avgPts)
	}
	if estimate.expectedPoints != 3900 {
		t.Errorf("estimate.expectedPoints = %v, want 3900", estimate.expectedPoints)
	}
	wantWall := []tile.Tile{
		tile.MustTileFromCode("1m"),
		tile.MustTileFromCode("5m"),
	}
	if !reflect.DeepEqual(wall, wantWall) {
		t.Errorf("wall = %v, want unchanged %v", wall, wantWall)
	}
}

func TestWinEstimatesFromShuffledWallRejectsInvalidInputs(t *testing.T) {
	candidates := []actionCandidate{{traceKey: "-1.1m"}}
	goalsByKey := map[string][]winEstimateGoal{"-1.1m": {}}
	wall := []tile.Tile{tile.MustTileFromCode("1m")}

	if _, err := winEstimatesFromShuffledWall(candidates, goalsByKey, wall, 1, 0, rand.New(rand.NewPCG(1, 0))); err == nil {
		t.Fatal("winEstimatesFromShuffledWall() accepted zero numTries")
	}
	if _, err := winEstimatesFromShuffledWall(candidates, goalsByKey, wall, 1, 1, nil); err == nil {
		t.Fatal("winEstimatesFromShuffledWall() accepted nil rng")
	}
	if _, err := winEstimatesFromShuffledWall(candidates, goalsByKey, wall, 2, 1, rand.New(rand.NewPCG(1, 0))); err == nil {
		t.Fatal("winEstimatesFromShuffledWall() accepted too many draws")
	}
}

func TestWinEstimatesFromState(t *testing.T) {
	candidates := []actionCandidate{
		{traceKey: "-1.1m"},
	}
	goalsByKey := map[string][]winEstimateGoal{
		"-1.1m": {
			{
				Goal: service.Goal{
					RequiredVector: hand.TileCounts34{0: 1},
				},
				points: 1000,
			},
		},
	}
	visibleTiles := make([]tile.Tile, 0, 135)
	for id := range tile.NumTileType34 {
		count := 4
		if id == 0 {
			count = 3
		}
		for range count {
			visibleTiles = append(visibleTiles, tile.MustTileFromID(id))
		}
	}
	state := stubWinEstimateStateViewer{
		turn:         0,
		visibleTiles: visibleTiles,
	}
	turnDistribution := make([]float64, numTurnDistributionEntries)
	turnDistribution[0] = 1
	stats := stubManueStats{
		turnDistribution: turnDistribution,
	}

	got, err := winEstimatesFromState(
		stats,
		state,
		seat.MustSeat(0),
		candidates,
		goalsByKey,
		3,
		rand.New(rand.NewPCG(1, 0)),
	)
	if err != nil {
		t.Fatalf("winEstimatesFromState() failed: %v", err)
	}

	estimate := got["-1.1m"]
	if estimate.prob != 1 {
		t.Errorf("estimate.prob = %v, want 1", estimate.prob)
	}
	if estimate.avgPts != 1000 {
		t.Errorf("estimate.avgPts = %v, want 1000", estimate.avgPts)
	}
	if estimate.expectedPoints != 1000 {
		t.Errorf("estimate.expectedPoints = %v, want 1000", estimate.expectedPoints)
	}
}

func TestWinEstimatesFromStateReturnsErrorWithInvalidVisibleTiles(t *testing.T) {
	_, err := winEstimatesFromState(
		stubManueStats{turnDistribution: fullTurnDistribution(1)},
		stubWinEstimateStateViewer{visibleTiles: []tile.Tile{tile.MustTileFromCode("?")}},
		seat.MustSeat(0),
		[]actionCandidate{{traceKey: "-1.1m"}},
		map[string][]winEstimateGoal{"-1.1m": {}},
		1,
		rand.New(rand.NewPCG(1, 0)),
	)
	if err == nil {
		t.Fatal("winEstimatesFromState() succeeded unexpectedly")
	}
}
