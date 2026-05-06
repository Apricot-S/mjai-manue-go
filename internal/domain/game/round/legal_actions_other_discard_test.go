package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
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

func TestState_LegalActions_OnOtherDiscardIncludesRon(t *testing.T) {
	hands := newValidHands()
	hands[1] = ronWithTanyaoHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := *seat.MustSeat(0)
	actor := *seat.MustSeat(1)
	winningTile := *tile.MustTileFromCode("3p")
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if !containsWin(got, actor, target, "3p") {
		t.Error("LegalActions() does not contain Win, want ron with tanyao")
	}
}

func TestState_LegalActions_OnOtherDiscardExcludesRonWithoutYaku(t *testing.T) {
	hands := newValidHands()
	hands[1] = ronWithoutYakuHandForLegalActionsTest()
	s := mustNewRoundStateForTest(t, hands)
	target := *seat.MustSeat(0)
	actor := *seat.MustSeat(1)
	winningTile := *tile.MustTileFromCode("9s")
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got, err := s.LegalActions(actor)
	if err != nil {
		t.Fatalf("LegalActions() failed: %v", err)
	}
	if containsWin(got, actor, target, "9s") {
		t.Error("LegalActions() contains Win, want no-yaku ron excluded")
	}
}

func ronWithTanyaoHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("2m"),
		*tile.MustTileFromCode("3m"),
		*tile.MustTileFromCode("4m"),
		*tile.MustTileFromCode("2p"),
		*tile.MustTileFromCode("4p"),
		*tile.MustTileFromCode("3s"),
		*tile.MustTileFromCode("4s"),
		*tile.MustTileFromCode("5s"),
		*tile.MustTileFromCode("6s"),
		*tile.MustTileFromCode("7s"),
		*tile.MustTileFromCode("8s"),
		*tile.MustTileFromCode("6m"),
		*tile.MustTileFromCode("6m"),
	}
}

func ronWithoutYakuHandForLegalActionsTest() [common.InitHandSize]tile.Tile {
	return [common.InitHandSize]tile.Tile{
		*tile.MustTileFromCode("1m"),
		*tile.MustTileFromCode("1m"),
		*tile.MustTileFromCode("1m"),
		*tile.MustTileFromCode("2p"),
		*tile.MustTileFromCode("3p"),
		*tile.MustTileFromCode("4p"),
		*tile.MustTileFromCode("3s"),
		*tile.MustTileFromCode("4s"),
		*tile.MustTileFromCode("5s"),
		*tile.MustTileFromCode("6s"),
		*tile.MustTileFromCode("6s"),
		*tile.MustTileFromCode("6s"),
		*tile.MustTileFromCode("9s"),
	}
}
