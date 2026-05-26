package ai

import (
	"strings"
	"testing"

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
