package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_LegalActions_OnOtherDiscardNoAction(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	target := *seat.MustSeat(0)
	actor := *seat.MustSeat(2)
	discardedTile := *tile.MustTileFromCode("6m")
	if err := s.Apply(event.NewDraw(target, discardedTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, discardedTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("LegalActions() = %v, want empty", got)
	}
}
