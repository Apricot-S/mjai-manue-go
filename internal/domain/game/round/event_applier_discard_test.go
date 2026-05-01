package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Discard(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	discardedTile := *tile.MustTileFromCode("6m")

	if err := s.Apply(event.NewDraw(actor, discardedTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	before := s.NumLeftTiles()

	if err := s.Apply(event.NewDiscard(actor, discardedTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if got := s.NumLeftTiles(); got != before {
		t.Errorf("NumLeftTiles() = %d, want %d", got, before)
	}
	if got := s.Player(actor).DrawnTile(); got != nil {
		t.Fatalf("DrawnTile() = %v, want nil", got)
	}
	if s.Player(actor).CanDiscard() {
		t.Error("CanDiscard() = true, want false")
	}
	if got := s.Player(actor).River(); len(got) != 1 || got[0].ID() != discardedTile.ID() {
		t.Fatalf("River() = %v, want [%v]", got, discardedTile)
	}
	if got := s.Player(actor).DiscardedTiles(); len(got) != 1 || got[0].ID() != discardedTile.ID() {
		t.Fatalf("DiscardedTiles() = %v, want [%v]", got, discardedTile)
	}
}
