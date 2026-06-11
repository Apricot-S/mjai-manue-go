package ai

import (
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestDangerSceneEvaluateReturnsErrorWithUnknownFeature(t *testing.T) {
	_, err := (dangerScene{}).evaluate("unknown_feature", tile.MustTileFromCode("5m"))
	if err == nil {
		t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "unknown_feature") {
		t.Errorf("dangerScene.evaluate() error = %v, want feature name", err)
	}
}

func TestDangerSceneEvaluateReturnsErrorWithInvalidFeatureInteger(t *testing.T) {
	_, err := (dangerScene{}).evaluate("visible>=invalid", tile.MustTileFromCode("5m"))
	if err == nil {
		t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "visible>=invalid") {
		t.Errorf("dangerScene.evaluate() error = %v, want feature name", err)
	}
}

func TestDangerSceneEvaluateRejectsOuterPreRiichiFeatureWithTrailingText(t *testing.T) {
	_, err := (dangerScene{}).evaluate("1_outer_prereach_sutehai_invalid", tile.MustTileFromCode("5m"))
	if err == nil {
		t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
	}
}

func TestDangerSceneEvaluateRejectsInnerPreRiichiFeatureWithTrailingText(t *testing.T) {
	_, err := (dangerScene{}).evaluate("1_inner_prereach_sutehai_invalid", tile.MustTileFromCode("5m"))
	if err == nil {
		t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
	}
}

func TestDangerSceneEvaluateReturnsErrorWithInvalidOuterPreRiichiInteger(t *testing.T) {
	_, err := (dangerScene{}).evaluate("invalid_outer_prereach_sutehai", tile.MustTileFromCode("5m"))
	if err == nil {
		t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "invalid_outer_prereach_sutehai") {
		t.Errorf("dangerScene.evaluate() error = %v, want feature name", err)
	}
}

func TestDangerSceneEvaluateReturnsErrorWithInvalidInnerPreRiichiInteger(t *testing.T) {
	_, err := (dangerScene{}).evaluate("invalid_inner_prereach_sutehai", tile.MustTileFromCode("5m"))
	if err == nil {
		t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
	}
	if !strings.Contains(err.Error(), "invalid_inner_prereach_sutehai") {
		t.Errorf("dangerScene.evaluate() error = %v, want feature name", err)
	}
}

func TestDangerSceneEvaluateReturnsErrorWithInvalidNumberRangeInteger(t *testing.T) {
	for _, feature := range []string{"invalid<=n<=6", "4<=n<=invalid"} {
		t.Run(feature, func(t *testing.T) {
			_, err := (dangerScene{}).evaluate(feature, tile.MustTileFromCode("5m"))
			if err == nil {
				t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
			}
			if !strings.Contains(err.Error(), feature) {
				t.Errorf("dangerScene.evaluate() error = %v, want feature name", err)
			}
		})
	}
}

func TestDangerSceneEvaluateReturnsErrorWithInvalidNeighborPreRiichiInteger(t *testing.T) {
	for _, feature := range []string{
		"+-invalid_in_prereach_sutehais>=1",
		"+-1_in_prereach_sutehais>=invalid",
	} {
		t.Run(feature, func(t *testing.T) {
			_, err := (dangerScene{}).evaluate(feature, tile.MustTileFromCode("5m"))
			if err == nil {
				t.Fatal("dangerScene.evaluate() succeeded unexpectedly")
			}
			if !strings.Contains(err.Error(), feature) {
				t.Errorf("dangerScene.evaluate() error = %v, want feature name", err)
			}
		})
	}
}

func TestDangerSceneEvaluateKnownFeature(t *testing.T) {
	got, err := (dangerScene{}).evaluate("sangenpai", tile.MustTileFromCode("P"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate() failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate() = false, want true")
	}
}

func TestDangerSceneEvaluateFanpaiFeatures(t *testing.T) {
	scene := dangerScene{
		roundWind:  wind.East,
		targetWind: wind.East,
	}
	got, err := scene.evaluate("ryenfonpai", tile.MustTileFromCode("E"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(ryenfonpai) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(ryenfonpai) = false, want true")
	}

	got, err = scene.evaluate("fanpai", tile.MustTileFromCode("P"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(fanpai) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(fanpai) = false, want true")
	}
}

func TestDangerSceneEvaluateVisibilityFeatures(t *testing.T) {
	scene := dangerScene{
		visibleTiles: []tile.Tile{
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5mr"),
			tile.MustTileFromCode("2m"),
		},
	}

	got, err := scene.evaluate("visible>=1", tile.MustTileFromCode("5m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(visible>=1) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(visible>=1) = false, want true")
	}

	got, err = scene.evaluate("suji_visible<=0", tile.MustTileFromCode("5m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(suji_visible<=0) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(suji_visible<=0) = false, want true because 8m is not visible")
	}
}

func TestDangerSceneEvaluateSujiInTehaisCountsEachSujiSeparately(t *testing.T) {
	tests := []struct {
		name     string
		discard  string
		selfHand []string
		want     bool
	}{
		{
			name:     "combined suji count does not satisfy threshold",
			discard:  "5m",
			selfHand: []string{"2m", "8m"},
			want:     false,
		},
		{
			name:     "one suji satisfies threshold",
			discard:  "5m",
			selfHand: []string{"2m", "2m"},
			want:     true,
		},
		{
			name:     "honor has no suji",
			discard:  "E",
			selfHand: []string{"E", "E"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selfHand := make([]tile.Tile, 0, len(tt.selfHand))
			for _, code := range tt.selfHand {
				selfHand = append(selfHand, tile.MustTileFromCode(code))
			}
			scene := dangerScene{selfHand: selfHand}

			got, err := scene.evaluate("suji_in_tehais>=2", tile.MustTileFromCode(tt.discard))
			if err != nil {
				t.Fatalf("dangerScene.evaluate(suji_in_tehais>=2) failed: %v", err)
			}
			if got != tt.want {
				t.Errorf("dangerScene.evaluate(suji_in_tehais>=2) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDangerSceneEvaluateChanceFeatureUsesKabeTiles(t *testing.T) {
	scene := dangerScene{
		visibleTiles: []tile.Tile{
			tile.MustTileFromCode("2m"),
			tile.MustTileFromCode("2m"),
			tile.MustTileFromCode("2m"),
			tile.MustTileFromCode("2m"),
		},
	}

	got, err := scene.evaluate("chances<=0", tile.MustTileFromCode("1m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(chances<=0) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(chances<=0) = false, want true with 2m kabe")
	}

	got, err = scene.evaluate("chances<=0", tile.MustTileFromCode("5m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(chances<=0 middle) failed: %v", err)
	}
	if got {
		t.Error("dangerScene.evaluate(chances<=0 middle) = true, want false for 4-6")
	}
}

func TestDangerSceneEvaluateAida4KenMatchesOriginal(t *testing.T) {
	scene := dangerScene{
		preRiichiTiles: []tile.Tile{
			tile.MustTileFromCode("1p"),
			tile.MustTileFromCode("6p"),
		},
	}

	wants := map[string]bool{
		"1p": false,
		"2p": true,
		"3p": false,
		"4p": false,
		"5p": true,
		"6p": false,
		"7p": false,
		"8p": false,
		"9p": false,
		"2m": false,
	}
	for code, want := range wants {
		t.Run(code, func(t *testing.T) {
			got, err := scene.evaluate("aida4ken", tile.MustTileFromCode(code))
			if err != nil {
				t.Fatalf("dangerScene.evaluate(aida4ken) failed: %v", err)
			}
			if got != want {
				t.Errorf("dangerScene.evaluate(aida4ken) = %v, want %v", got, want)
			}
		})
	}
}

func TestDangerSceneEvaluateOuterPreRiichiMatchesOriginalDirection(t *testing.T) {
	scene := dangerScene{
		preRiichiTiles: []tile.Tile{tile.MustTileFromCode("4m")},
	}

	got, err := scene.evaluate("1_outer_prereach_sutehai", tile.MustTileFromCode("3m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(1_outer_prereach_sutehai) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(1_outer_prereach_sutehai) = false, want true for original direction")
	}

	got, err = scene.evaluate("1_inner_prereach_sutehai", tile.MustTileFromCode("3m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(1_inner_prereach_sutehai) failed: %v", err)
	}
	if got {
		t.Error("dangerScene.evaluate(1_inner_prereach_sutehai) = true, want false without 2m")
	}
}

func TestDangerSceneEvaluateSameTypeInPreRiichiCountsDistinctSuitNumbers(t *testing.T) {
	scene := dangerScene{
		preRiichiTiles: []tile.Tile{
			tile.MustTileFromCode("5m"),
			tile.MustTileFromCode("5mr"),
			tile.MustTileFromCode("7m"),
			tile.MustTileFromCode("1p"),
		},
	}

	got, err := scene.evaluate("same_type_in_prereach>=3", tile.MustTileFromCode("2m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(same_type_in_prereach>=3) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(same_type_in_prereach>=3) = false, want true for two distinct prereach manzu plus discard")
	}

	got, err = scene.evaluate("same_type_in_prereach>=4", tile.MustTileFromCode("2m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(same_type_in_prereach>=4) failed: %v", err)
	}
	if got {
		t.Error("dangerScene.evaluate(same_type_in_prereach>=4) = true, want false because duplicate 5m/5mr counts once")
	}

	got, err = scene.evaluate("same_type_in_prereach>=1", tile.MustTileFromCode("E"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(same_type_in_prereach>=1 honors) failed: %v", err)
	}
	if got {
		t.Error("dangerScene.evaluate(same_type_in_prereach>=1) = true, want false for honors")
	}
}

func TestDangerSceneEvaluateNeighborPreRiichiMatchesOriginalRange(t *testing.T) {
	tests := []struct {
		name           string
		feature        string
		discard        string
		preRiichiTiles []string
		want           bool
	}{
		{
			name:           "includes discard itself",
			feature:        "+-1_in_prereach_sutehais>=1",
			discard:        "5m",
			preRiichiTiles: []string{"5m"},
			want:           true,
		},
		{
			name:           "distance two includes intermediate numbers",
			feature:        "+-2_in_prereach_sutehais>=2",
			discard:        "5m",
			preRiichiTiles: []string{"4m", "6m"},
			want:           true,
		},
		{
			name:           "counts duplicate tiles as one number",
			feature:        "+-1_in_prereach_sutehais>=2",
			discard:        "1p",
			preRiichiTiles: []string{"2p", "2p"},
			want:           false,
		},
		{
			name:           "counts distinct numbers",
			feature:        "+-1_in_prereach_sutehais>=2",
			discard:        "2p",
			preRiichiTiles: []string{"1p", "3p"},
			want:           true,
		},
		{
			name:           "excludes tiles outside bounded range",
			feature:        "+-2_in_prereach_sutehais>=2",
			discard:        "1s",
			preRiichiTiles: []string{"3s", "4s"},
			want:           false,
		},
		{
			name:           "honor has no numbered neighbors",
			feature:        "+-1_in_prereach_sutehais>=1",
			discard:        "E",
			preRiichiTiles: []string{"E"},
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preRiichiTiles := make([]tile.Tile, 0, len(tt.preRiichiTiles))
			for _, code := range tt.preRiichiTiles {
				preRiichiTiles = append(preRiichiTiles, tile.MustTileFromCode(code))
			}
			scene := dangerScene{preRiichiTiles: preRiichiTiles}

			got, err := scene.evaluate(tt.feature, tile.MustTileFromCode(tt.discard))
			if err != nil {
				t.Fatalf("dangerScene.evaluate(%s) failed: %v", tt.feature, err)
			}
			if got != tt.want {
				t.Errorf("dangerScene.evaluate(%s) = %v, want %v", tt.feature, got, tt.want)
			}
		})
	}
}

func TestNewDangerSceneKeepsPreRiichiTilesEmptyWithoutRiichi(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(1)
	var players [common.NumPlayers]player.PlayerViewer
	players[self.Index()] = stubPlayerViewer{}
	players[target.Index()] = stubPlayerViewer{
		discardedTiles: []tile.Tile{tile.MustTileFromCode("1m")},
	}
	state := stubCandidateEvaluationStateViewer{
		players: players,
	}

	scene := newDangerScene(state, self, target)
	got, err := scene.evaluate("prereach_suji", tile.MustTileFromCode("4m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(prereach_suji) failed: %v", err)
	}
	if got {
		t.Error("dangerScene.evaluate(prereach_suji) = true, want false before target riichi")
	}
}

func TestNewDangerSceneUsesTilesThroughRiichiDiscard(t *testing.T) {
	self := seat.MustSeat(0)
	target := seat.MustSeat(1)
	var players [common.NumPlayers]player.PlayerViewer
	players[self.Index()] = stubPlayerViewer{}
	players[target.Index()] = stubPlayerViewer{
		discardedTiles: []tile.Tile{
			tile.MustTileFromCode("1m"),
			tile.MustTileFromCode("7m"),
			tile.MustTileFromCode("2p"),
			tile.MustTileFromCode("8p"),
		},
		riichiDiscardedTilesIndex: 1,
		hasRiichiDiscardIndex:     true,
	}
	state := stubCandidateEvaluationStateViewer{
		players: players,
	}

	scene := newDangerScene(state, self, target)
	got, err := scene.evaluate("prereach_suji", tile.MustTileFromCode("4m"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(prereach_suji) failed: %v", err)
	}
	if !got {
		t.Error("dangerScene.evaluate(prereach_suji) = false, want true for riichi discard suji")
	}

	got, err = scene.evaluate("prereach_suji", tile.MustTileFromCode("5p"))
	if err != nil {
		t.Fatalf("dangerScene.evaluate(prereach_suji late discard) failed: %v", err)
	}
	if got {
		t.Error("dangerScene.evaluate(prereach_suji) = true, want false for discard after riichi declaration")
	}
}

func TestDecisionTreeDangerEstimator_SafeTileSkipsSceneBuild(t *testing.T) {
	self := seat.MustSeat(0)
	winner := seat.MustSeat(1)
	discard := tile.MustTileFromCode("5mr")
	state := safeOnlyStateViewer{
		winner:   winner,
		safeTile: discard.RemoveRed(),
	}
	estimator := NewDangerEstimator(stubDangerTreeFeature{
		feature: "unknown_feature",
		negative: stubDangerTreeLeaf{
			prob: 0.25,
		},
		positive: stubDangerTreeLeaf{
			prob: 0.75,
		},
	})

	got, err := estimator.EstimateDealInProb(state, self, winner, discard)
	if err != nil {
		t.Fatalf("EstimateDealInProb() failed: %v", err)
	}
	if got != 0 {
		t.Errorf("EstimateDealInProb() = %v, want 0 for safe tile", got)
	}
}

type safeOnlyStateViewer struct {
	round.StateViewer
	winner   seat.Seat
	safeTile tile.Tile
}

func (s safeOnlyStateViewer) SafeTiles(playerSeat seat.Seat) tile.Tiles {
	if playerSeat == s.winner {
		return tile.Tiles{s.safeTile}
	}
	return nil
}

type stubDangerTreeFeature struct {
	feature  string
	negative DangerTreeNode
	positive DangerTreeNode
}

func (s stubDangerTreeFeature) LeafProb() (float64, bool) {
	return 0, false
}

func (s stubDangerTreeFeature) Feature() (string, bool) {
	return s.feature, true
}

func (s stubDangerTreeFeature) NegativeNode() DangerTreeNode {
	return s.negative
}

func (s stubDangerTreeFeature) PositiveNode() DangerTreeNode {
	return s.positive
}
