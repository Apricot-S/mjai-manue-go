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

func TestDangerSceneEvaluateOuterPrereachMatchesOriginalDirection(t *testing.T) {
	scene := dangerScene{
		prereachTiles: []tile.Tile{tile.MustTileFromCode("4m")},
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

func TestNewDangerSceneKeepsPrereachTilesEmptyWithoutRiichi(t *testing.T) {
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
