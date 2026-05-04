package round

import (
	"fmt"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
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

func TestState_LegalActions_NotPendingActor(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	if err := s.Apply(event.NewDraw(*seat.MustSeat(0), *tile.MustTileFromCode("6m"))); err != nil {
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

func TestState_LegalActions_PendingDiscard(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	drawnTile := *tile.MustTileFromCode("6m")
	if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}

	want := map[string]bool{
		"1m:false": false,
		"1p:false": false,
		"1s:false": false,
		"2m:false": false,
		"2p:false": false,
		"2s:false": false,
		"3m:false": false,
		"3p:false": false,
		"3s:false": false,
		"4m:false": false,
		"4p:false": false,
		"4s:false": false,
		"5m:false": false,
		"6m:true":  false,
	}
	if len(got) != len(want) {
		t.Fatalf("LegalActions() length = %d, want %d: %v", len(got), len(want), got)
	}
	for _, a := range got {
		discard, ok := a.(*action.Discard)
		if !ok {
			t.Fatalf("LegalActions() contains %T, want only *action.Discard", a)
		}
		if discard.Actor() != actor {
			t.Errorf("Discard.Actor() = %v, want %v", discard.Actor(), actor)
		}
		key := fmt.Sprintf("%s:%t", discard.Tile(), discard.Tsumogiri())
		if _, ok := want[key]; !ok {
			t.Errorf("unexpected discard action: %s", key)
			continue
		}
		want[key] = true
	}
	for key, found := range want {
		if !found {
			t.Errorf("missing discard action: %s", key)
		}
	}
}

func TestState_LegalActions_AfterRiichiAcceptedAllowsOnlyTsumogiri(t *testing.T) {
	hands := newValidHands()
	hands[0] = riichiReadyHandForTest()
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(0)
	firstDraw := *tile.MustTileFromCode("S")
	firstDiscard := *tile.MustTileFromCode("W")
	if err := s.Apply(event.NewDraw(actor, firstDraw)); err != nil {
		t.Fatalf("Apply(first Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichi(actor)); err != nil {
		t.Fatalf("Apply(Riichi) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, firstDiscard, false)); err != nil {
		t.Fatalf("Apply(first Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewRiichiAccepted(actor, nil, nil)); err != nil {
		t.Fatalf("Apply(RiichiAccepted) failed: %v", err)
	}

	for i := 1; i < 4; i++ {
		other := *seat.MustSeat(i)
		drawnTile := *tile.MustTileFromCode("6m")
		if err := s.Apply(event.NewDraw(other, drawnTile)); err != nil {
			t.Fatalf("Apply(other Draw %d) failed: %v", i, err)
		}
		if err := s.Apply(event.NewDiscard(other, drawnTile, true)); err != nil {
			t.Fatalf("Apply(other Discard %d) failed: %v", i, err)
		}
	}

	secondDraw := *tile.MustTileFromCode("7m")
	if err := s.Apply(event.NewDraw(actor, secondDraw)); err != nil {
		t.Fatalf("Apply(second Draw) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("LegalActions() length = %d, want 1: %v", len(got), got)
	}
	discard, ok := got[0].(*action.Discard)
	if !ok {
		t.Fatalf("LegalActions()[0] = %T, want *action.Discard", got[0])
	}
	if discard.Tile() != secondDraw {
		t.Errorf("Discard.Tile() = %v, want %v", discard.Tile(), secondDraw)
	}
	if !discard.Tsumogiri() {
		t.Error("Discard.Tsumogiri() = false, want true")
	}
}

func TestState_LegalActions_InvalidatesCacheAfterApply(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	drawnTile := *tile.MustTileFromCode("6m")
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
	if err := s.Apply(event.NewDraw(actor, *tile.MustTileFromCode("6m"))); err != nil {
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
