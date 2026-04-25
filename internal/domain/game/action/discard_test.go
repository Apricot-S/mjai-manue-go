package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewDiscard(t *testing.T) {
	actor := *seat.MustSeat(1)
	discardedTile := *tile.MustTileFromCode("5mr")

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
	actor := *seat.MustSeat(1)
	unknownTile := *tile.MustTileFromCode("?")

	if _, err := action.NewDiscard(actor, unknownTile, true); err == nil {
		t.Fatal("NewDiscard() succeeded unexpectedly")
	}
}
