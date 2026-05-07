package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestState_Apply_Chii(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(1)
	target := *seat.MustSeat(0)
	taken := tile.MustTileFromCode("4m")
	consumed := [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewChii(actor, target, taken, consumed)); err != nil {
		t.Fatalf("Apply(Chii) failed: %v", err)
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
	if !s.Player(actor).CanDiscard() {
		t.Error("actor CanDiscard() = false, want true after Chii")
	}
	if s.Player(actor).IsConcealed() {
		t.Error("actor IsConcealed() = true, want false after Chii")
	}
}

func TestState_Apply_Chii_ReturnsErrorWhenActorAndTargetAreSame(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	taken := tile.MustTileFromCode("4m")
	consumed := [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")}

	if err := s.Apply(event.NewDraw(actor, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewChii(actor, actor, taken, consumed)); err == nil {
		t.Fatal("Apply(Chii) succeeded unexpectedly")
	}

	if got := s.Player(actor).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("actor River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}

func TestState_Apply_Chii_ReturnsErrorForLastDiscard(t *testing.T) {
	players := newVisiblePlayersForTest(t, newValidHands())
	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		*seat.MustSeat(0),
		*seat.MustSeat(0),
		tile.Tiles{tile.MustTileFromCode("E")},
		1,
		players,
	)
	actor := *seat.MustSeat(1)
	target := *seat.MustSeat(0)
	taken := tile.MustTileFromCode("4m")
	consumed := [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if got := s.NumLeftTiles(); got != 0 {
		t.Fatalf("NumLeftTiles() = %d, want 0", got)
	}

	if err := s.Apply(event.NewChii(actor, target, taken, consumed)); err == nil {
		t.Fatal("Apply(Chii) succeeded unexpectedly")
	}

	if got := s.Player(target).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("target River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}

func TestState_Apply_Chii_ReturnsErrorWhenActorIsNotShimochaOfTarget(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(2)
	target := *seat.MustSeat(0)
	taken := tile.MustTileFromCode("4m")
	consumed := [2]tile.Tile{tile.MustTileFromCode("2m"), tile.MustTileFromCode("3m")}

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewChii(actor, target, taken, consumed)); err == nil {
		t.Fatal("Apply(Chii) succeeded unexpectedly")
	}

	if got := s.Player(target).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("target River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}
