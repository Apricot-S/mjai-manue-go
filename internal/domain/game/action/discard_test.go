package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewDiscard(t *testing.T) {
	actor := seat.MustSeat(1)
	discardedTile := tile.MustTileFromCode("5mr")

	got, err := action.NewDiscard(actor, discardedTile, true)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	if got.Actor() != actor {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
	if got.Tile() != discardedTile {
		t.Errorf("Tile() = %v, want %v", got.Tile(), discardedTile)
	}
	if !got.Tsumogiri() {
		t.Error("Tsumogiri() = false, want true")
	}
}

func TestNewDiscard_UnknownTile(t *testing.T) {
	actor := seat.MustSeat(1)
	unknownTile := tile.MustTileFromCode("?")

	tests := []struct {
		name          string
		discardedTile tile.Tile
		tsumogiri     bool
	}{
		{
			name:          "discarded tile",
			discardedTile: unknownTile,
			tsumogiri:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := action.NewDiscard(actor, tt.discardedTile, tt.tsumogiri); err == nil {
				t.Error("NewDiscard() succeeded unexpectedly")
			}
		})
	}
}
