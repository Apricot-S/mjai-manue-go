package main

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func mustTiles(codes ...string) tile.Tiles {
	tiles := make(tile.Tiles, len(codes))
	for i, code := range codes {
		tiles[i] = tile.MustTileFromCode(code)
	}
	return tiles
}

func TestSceneEvaluateRubyOnlyUrasujiOf5(t *testing.T) {
	scene := NewSceneFromParams(nil, mustTiles("1p"), nil, nil, mustTiles("1p", "5s"), wind.East, wind.South)

	got, err := scene.evaluate("urasuji_of_5", tile.MustTileFromCode("1s"))
	if err != nil {
		t.Fatalf("evaluate() error = %v", err)
	}
	if !got {
		t.Errorf("evaluate(urasuji_of_5, 1s) = false, want true")
	}
}

func TestSceneEvaluateUrasujiOf5DoesNotMutatePreRiichiTiles(t *testing.T) {
	scene := NewSceneFromParams(nil, nil, nil, nil, mustTiles("1p", "5s"), wind.East, wind.South)

	if _, err := scene.evaluate("urasuji_of_5", tile.MustTileFromCode("1s")); err != nil {
		t.Fatalf("evaluate(urasuji_of_5) error = %v", err)
	}
	got, err := scene.evaluate("urasuji", tile.MustTileFromCode("2p"))
	if err != nil {
		t.Fatalf("evaluate(urasuji) error = %v", err)
	}
	if !got {
		t.Errorf("evaluate(urasuji, 2p) after urasuji_of_5 = false, want true")
	}
}

func TestSceneEvaluateSameTypeInPrereachMatchesRubyTool(t *testing.T) {
	scene := NewSceneFromParams(nil, nil, nil, nil, mustTiles("1p", "3p"), wind.East, wind.South)

	got, err := scene.evaluate("same_type_in_prereach>=3", tile.MustTileFromCode("5p"))
	if err != nil {
		t.Fatalf("evaluate() error = %v", err)
	}
	if got {
		t.Errorf("evaluate(same_type_in_prereach>=3) = true, want false")
	}
}

func TestSceneEvaluateInnerPreRiichiDiscardMatchesRubyTool(t *testing.T) {
	scene := NewSceneFromParams(nil, nil, nil, nil, mustTiles("6p", "7p"), wind.East, wind.South)

	tests := []struct {
		name    string
		feature string
		discard tile.Tile
		want    bool
	}{
		{
			name:    "one inner from five",
			feature: "1_inner_prereach_sutehai",
			discard: tile.MustTileFromCode("5p"),
			want:    true,
		},
		{
			name:    "two inner from five",
			feature: "2_inner_prereach_sutehai",
			discard: tile.MustTileFromCode("5p"),
			want:    true,
		},
		{
			name:    "two outer does not match center",
			feature: "2_outer_prereach_sutehai",
			discard: tile.MustTileFromCode("5p"),
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := scene.evaluate(tt.feature, tt.discard)
			if err != nil {
				t.Fatalf("evaluate() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("evaluate(%s, %s) = %v, want %v", tt.feature, tt.discard, got, tt.want)
			}
		})
	}
}

func TestSceneFeatureVectorTreatsRedFiveAsNormalFive(t *testing.T) {
	scene := NewSceneFromParams(mustTiles("5mr"), nil, nil, nil, nil, wind.East, wind.South)

	redVector, err := scene.FeatureVector(tile.MustTileFromCode("5mr"))
	if err != nil {
		t.Fatalf("FeatureVector() error = %v", err)
	}
	normalVector, err := scene.FeatureVector(tile.MustTileFromCode("5m"))
	if err != nil {
		t.Fatalf("FeatureVector() error = %v", err)
	}
	if redVector.Cmp(normalVector) != 0 {
		t.Errorf("FeatureVector(5mr) = %v, want %v", redVector, normalVector)
	}
	if !GetFeatureValue(redVector, "5<=n<=5") {
		t.Errorf("FeatureVector(5mr) did not set 5<=n<=5")
	}
}

func TestSceneCandidatesTreatRedFiveAsNormalFive(t *testing.T) {
	scene := NewSceneFromParams(mustTiles("5mr", "5m"), nil, nil, nil, nil, wind.East, wind.South)

	candidates := scene.Candidates()
	if len(candidates) != 1 {
		t.Fatalf("len(Candidates()) = %d, want 1: %v", len(candidates), candidates)
	}
	if candidates[0] != tile.MustTileFromCode("5m") {
		t.Errorf("Candidates()[0] = %s, want 5m", candidates[0])
	}
}
