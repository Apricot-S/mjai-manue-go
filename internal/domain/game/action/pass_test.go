package action_test

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func TestNewPass(t *testing.T) {
	actor := seat.MustSeat(2)

	got := action.NewPass(actor)

	if got.Actor() != actor {
		t.Errorf("Actor() = %v, want %v", got.Actor(), actor)
	}
}
