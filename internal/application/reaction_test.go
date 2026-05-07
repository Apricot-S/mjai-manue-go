package application_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/application"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestNewNoReaction(t *testing.T) {
	got := application.NewNoReaction()
	if got.Kind() != application.ReactionNone {
		t.Errorf("Kind() = %v, want %v", got.Kind(), application.ReactionNone)
	}
	if got.Action() != nil {
		t.Errorf("Action() = %v, want nil", got.Action())
	}
}

func TestNewActionReaction(t *testing.T) {
	discard, err := action.NewDiscard(seat.MustSeat(0), tile.MustTileFromCode("1m"), true)
	if err != nil {
		t.Fatalf("action.NewDiscard() failed: %v", err)
	}

	got := application.NewActionReaction(discard, "log")
	if got.Kind() != application.ReactionAction {
		t.Errorf("Kind() = %v, want %v", got.Kind(), application.ReactionAction)
	}
	if got.Action() != discard {
		t.Errorf("Action() = %v, want %v", got.Action(), discard)
	}
	if got.Log() != "log" {
		t.Errorf("Log() = %q, want log", got.Log())
	}
}
