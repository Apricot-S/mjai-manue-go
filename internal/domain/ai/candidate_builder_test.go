package ai

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestCandidateShantenUsesThrowableVector(t *testing.T) {
	got := candidateShanten(
		tile.MustTileFromCode("1m"),
		3,
		[]service.Goal{
			{Shanten: 2, ThrowableVector: hand.TileCounts34{0: 1}},
			{Shanten: 1, ThrowableVector: hand.TileCounts34{0: 1}},
			{Shanten: 0, ThrowableVector: hand.TileCounts34{1: 1}},
		},
	)
	if got != 1 {
		t.Errorf("candidateShanten() = %d, want 1", got)
	}
}

func TestCandidateShantenReturnsBaseForNone(t *testing.T) {
	got := candidateShanten(tile.MustTileFromCode("?"), 3, nil)
	if got != 3 {
		t.Errorf("candidateShanten(?) = %d, want 3", got)
	}
}

func TestCandidateShantenReturnsInfinityWhenTileIsNotThrowable(t *testing.T) {
	got := candidateShanten(
		tile.MustTileFromCode("1m"),
		3,
		[]service.Goal{{Shanten: 0, ThrowableVector: hand.TileCounts34{1: 1}}},
	)
	if got != service.InfinityShanten {
		t.Errorf("candidateShanten() = %d, want InfinityShanten", got)
	}
}
