package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewWin(t *testing.T) {
	actor := seat.MustSeat(1)
	target := seat.MustSeat(0)
	winningTile := tile.MustTileFromCode("5mr")

	got, err := action.NewWin(actor, target, winningTile)
	if err != nil {
		t.Fatalf("NewWin() failed: %v", err)
	}

	if got.Actor() != actor {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
	if got.Target() != target {
		t.Errorf("Target() = %v, want %v", got.Target(), target)
	}
	if got.WinningTile() != winningTile {
		t.Errorf("WinningTile() = %v, want %v", got.WinningTile(), winningTile)
	}
}

func TestNewWin_UnknownTile(t *testing.T) {
	actor := seat.MustSeat(1)
	target := seat.MustSeat(0)
	unknown := tile.MustTileFromCode("?")

	tests := []struct {
		name        string
		winningTile tile.Tile
	}{
		{
			name:        "winning tile",
			winningTile: unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := action.NewWin(actor, target, tt.winningTile); err == nil {
				t.Error("NewWin() succeeded unexpectedly")
			}
		})
	}
}
