package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func calledKanHandsForTest() [common.NumPlayers][common.InitHandSize]tile.Tile {
	hands := newValidHands()
	hands[3] = [common.InitHandSize]tile.Tile{
		tile.MustTileFromCode("1m"), tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m"),
		tile.MustTileFromCode("4p"), tile.MustTileFromCode("5p"), tile.MustTileFromCode("6p"),
		tile.MustTileFromCode("7s"), tile.MustTileFromCode("8s"), tile.MustTileFromCode("9s"),
		tile.MustTileFromCode("E"), tile.MustTileFromCode("E"), tile.MustTileFromCode("E"),
		tile.MustTileFromCode("W"),
	}
	return hands
}

func TestState_Apply_CalledKan(t *testing.T) {
	s := mustNewRoundStateForTest(t, calledKanHandsForTest())
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewCalledKan(actor, target, taken, consumed)); err != nil {
		t.Fatalf("Apply(CalledKan) failed: %v", err)
	}

	if got := s.Player(target).River(); len(got) != 0 {
		t.Errorf("target River() = %v, want empty", got)
	}
	if got := s.Player(target).DiscardedTiles(); len(got) != 1 || got[0] != taken {
		t.Errorf("target DiscardedTiles() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 1 {
		t.Fatalf("actor Melds() length = %d, want 1", len(got))
	}
	if s.Player(actor).CanDiscard() {
		t.Error("actor CanDiscard() = true, want false after CalledKan")
	}
	if s.Player(actor).IsConcealed() {
		t.Error("actor IsConcealed() = true, want false after CalledKan")
	}
}

func TestState_Apply_CalledKan_ReturnsErrorWhenActorAndTargetAreSame(t *testing.T) {
	s := mustNewRoundStateForTest(t, calledKanHandsForTest())
	actor := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}

	if err := s.Apply(event.NewDraw(actor, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewCalledKan(actor, actor, taken, consumed)); err == nil {
		t.Fatal("Apply(CalledKan) succeeded unexpectedly")
	}

	if got := s.Player(actor).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("actor River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}

func TestState_Apply_CalledKan_ReturnsErrorForLastDiscard(t *testing.T) {
	players := newVisiblePlayersForTest(t, calledKanHandsForTest())
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		seat.MustSeat(0),
		seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		1,
		players,
	)
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if got := s.NumLeftTiles(); got != 0 {
		t.Fatalf("NumLeftTiles() = %d, want 0", got)
	}

	if err := s.Apply(event.NewCalledKan(actor, target, taken, consumed)); err == nil {
		t.Fatal("Apply(CalledKan) succeeded unexpectedly")
	}

	if got := s.Player(target).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("target River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}

func TestState_Apply_CalledKan_ReturnsErrorOnFifthKan(t *testing.T) {
	s := newStateForTestWithNumKans(mustNewRoundStateForTest(t, calledKanHandsForTest()), maxNumKan)
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewCalledKan(actor, target, taken, consumed)); err == nil {
		t.Fatal("Apply(CalledKan) succeeded unexpectedly")
	}

	if got := s.Player(target).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("target River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}

func TestState_Apply_CalledKan_ReturnsErrorWhenDiscardFollowsCalledKan(t *testing.T) {
	s := mustNewRoundStateForTest(t, calledKanHandsForTest())
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewCalledKan(actor, target, taken, consumed)); err != nil {
		t.Fatalf("Apply(CalledKan) failed: %v", err)
	}

	if err := s.Apply(event.NewDiscard(actor, tile.MustTileFromCode("W"), false)); err == nil {
		t.Fatal("Apply(Discard) succeeded unexpectedly")
	}

	if got := s.Player(actor).River(); len(got) != 0 {
		t.Fatalf("actor River() = %v, want empty", got)
	}
}

func TestState_Apply_CalledKan_AllowsDoraAfterReplacementTileDraw(t *testing.T) {
	s := mustNewRoundStateForTest(t, calledKanHandsForTest())
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}
	replacementTile := tile.MustTileFromCode("W")
	doraIndicator := tile.MustTileFromCode("6p")

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewCalledKan(actor, target, taken, consumed)); err != nil {
		t.Fatalf("Apply(CalledKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(actor, replacementTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	if err := s.Apply(event.NewDora(doraIndicator)); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}

	if got := s.DoraIndicators(); len(got) != 2 || got[1] != doraIndicator {
		t.Fatalf("DoraIndicators() = %v, want appended %v", got, doraIndicator)
	}
}

func TestState_Apply_CalledKan_AllowsDiscardAfterDoraReveal(t *testing.T) {
	s := mustNewRoundStateForTest(t, calledKanHandsForTest())
	actor := seat.MustSeat(3)
	target := seat.MustSeat(0)
	taken := tile.MustTileFromCode("E")
	consumed := [3]tile.Tile{taken, taken, taken}
	replacementTile := tile.MustTileFromCode("W")

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewCalledKan(actor, target, taken, consumed)); err != nil {
		t.Fatalf("Apply(CalledKan) failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(actor, replacementTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDora(tile.MustTileFromCode("6p"))); err != nil {
		t.Fatalf("Apply(Dora) failed: %v", err)
	}

	if err := s.Apply(event.NewDiscard(actor, replacementTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if got := s.Player(actor).River(); len(got) != 1 || got[0] != replacementTile {
		t.Fatalf("actor River() = %v, want [%v]", got, replacementTile)
	}
}
