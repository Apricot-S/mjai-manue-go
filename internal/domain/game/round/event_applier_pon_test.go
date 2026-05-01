package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestState_Apply_Pon(t *testing.T) {
	hands := newValidHands()
	hands[3][1] = *tile.MustTileFromCode("1s")
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(3)
	target := *seat.MustSeat(0)
	taken := *tile.MustTileFromCode("1s")

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewPon(actor, target, taken, [2]tile.Tile{taken, taken})); err != nil {
		t.Fatalf("Apply(Pon) failed: %v", err)
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
		t.Error("actor CanDiscard() = false, want true after Pon")
	}
	if s.Player(actor).IsConcealed() {
		t.Error("actor IsConcealed() = true, want false after Pon")
	}
}

func TestState_Apply_Pon_ReturnsErrorWhenActorAndTargetAreSame(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	taken := *tile.MustTileFromCode("1s")

	if err := s.Apply(event.NewDraw(actor, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	if err := s.Apply(event.NewPon(actor, actor, taken, [2]tile.Tile{taken, taken})); err == nil {
		t.Fatal("Apply(Pon) succeeded unexpectedly")
	}

	if got := s.Player(actor).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("actor River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}

func TestState_Apply_Pon_ReturnsErrorForLastDiscard(t *testing.T) {
	hands := newValidHands()
	hands[3][1] = *tile.MustTileFromCode("1s")
	players := newVisiblePlayersForTest(t, hands)

	s := NewStateForTest(
		wind.East,
		1,
		0,
		0,
		[common.NumPlayers]int{25000, 25000, 25000, 25000},
		*seat.MustSeat(0),
		*seat.MustSeat(0),
		tile.Tiles{*tile.MustTileFromCode("E")},
		1,
		players,
	)
	actor := *seat.MustSeat(3)
	target := *seat.MustSeat(0)
	taken := *tile.MustTileFromCode("1s")

	if err := s.Apply(event.NewDraw(target, taken)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, taken, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if got := s.NumLeftTiles(); got != 0 {
		t.Fatalf("NumLeftTiles() = %d, want 0", got)
	}

	if err := s.Apply(event.NewPon(actor, target, taken, [2]tile.Tile{taken, taken})); err == nil {
		t.Fatal("Apply(Pon) succeeded unexpectedly")
	}

	if got := s.Player(target).River(); len(got) != 1 || got[0] != taken {
		t.Errorf("target River() = %v, want [%v]", got, taken)
	}
	if got := s.Player(actor).Melds(); len(got) != 0 {
		t.Fatalf("actor Melds() length = %d, want 0", len(got))
	}
}
