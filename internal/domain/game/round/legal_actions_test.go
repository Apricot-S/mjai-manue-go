package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_LegalActions_NoPendingAction(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())

	got, err := s.LegalActions(*seat.MustSeat(0))
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("LegalActions() = %v, want empty", got)
	}
}

func TestState_LegalActions_ReturnsErrorForInvisiblePlayerWithoutPendingAction(t *testing.T) {
	hands := newValidHands()
	hands[0] = unknownHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)

	if _, err := s.LegalActions(*seat.MustSeat(0)); err == nil {
		t.Fatal("LegalActions() succeeded unexpectedly")
	}
}

func TestState_LegalActions_ReturnsErrorForPendingInvisiblePlayer(t *testing.T) {
	hands := newValidHands()
	hands[0] = unknownHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	if _, err := s.LegalActions(actor); err == nil {
		t.Fatal("LegalActions() succeeded unexpectedly")
	}
}

func TestState_LegalActions_NotPendingActor(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	if err := s.Apply(event.NewDraw(*seat.MustSeat(0), tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(*seat.MustSeat(1))
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("LegalActions() = %v, want empty", got)
	}
}

func TestState_LegalActions_InvalidatesCacheAfterApply(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	drawnTile := tile.MustTileFromCode("6m")
	if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	before, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() before discard failed: %v", err)
	}
	if len(before) == 0 {
		t.Fatal("LegalActions() before discard is empty")
	}

	if err := s.Apply(event.NewDiscard(actor, drawnTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	after, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() after discard failed: %v", err)
	}
	if len(after) != 0 {
		t.Fatalf("LegalActions() after discard = %v, want empty", after)
	}
}

func TestState_LegalActions_ReturnsSliceCopy(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	if err := s.Apply(event.NewDraw(actor, tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	first, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() first call failed: %v", err)
	}
	if len(first) == 0 {
		t.Fatal("LegalActions() first call is empty")
	}
	first[0] = nil

	second, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() second call failed: %v", err)
	}
	if second[0] == nil {
		t.Fatal("LegalActions() returned cache slice directly")
	}
}
