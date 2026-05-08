package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Discard(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := seat.MustSeat(0)
	discardedTile := tile.MustTileFromCode("6m")

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

func TestState_Apply_Discard_ReturnsErrorWhenActorIsNotPendingDiscardPlayer(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	drawActor := seat.MustSeat(0)
	discardActor := seat.MustSeat(1)
	drawnTile := tile.MustTileFromCode("6m")

	if err := s.Apply(event.NewDraw(drawActor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	if err := s.Apply(event.NewDiscard(discardActor, drawnTile, true)); err == nil {
		t.Fatal("Apply(Discard) succeeded unexpectedly")
	}

	if got := s.Player(drawActor).DrawnTile(); got == nil || got.ID() != drawnTile.ID() {
		t.Fatalf("draw actor DrawnTile() = %v, want %v", got, drawnTile)
	}
	if got := s.Player(discardActor).River(); len(got) != 0 {
		t.Fatalf("discard actor River() = %v, want empty", got)
	}
}

func TestState_Apply_Discard_ReturnsErrorForSwapCallTileAfterPon(t *testing.T) {
	hands := newValidHands()
	hands[3] = [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("5pr"), tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"), tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("5p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("S"),
		tile.MustTileFromCode("W"),
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("5p")
	swapCallTile := tile.MustTileFromCode("5pr")

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewPon(actor, target, taken, [2]tile.Tile{taken, taken})); err != nil {
		t.Fatalf("Apply(Pon) failed: %v", err)
	}

	if err := s.Apply(event.NewDiscard(actor, swapCallTile, false)); err == nil {
		t.Fatal("Apply(Discard) succeeded unexpectedly")
	}

	if got := s.Player(actor).River(); len(got) != 0 {
		t.Fatalf("actor River() = %v, want empty", got)
	}
	if !s.Player(actor).CanDiscard() {
		t.Error("actor CanDiscard() = false, want true after failed Discard")
	}
}
