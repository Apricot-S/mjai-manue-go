package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func TestState_RenderBoard(t *testing.T) {
	hands := newValidHands()
	ev := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		*seat.MustSeat(0),
		*tile.MustTileFromCode("E"),
		&[4]int{25000, 25000, 25000, 25000},
		hands,
	)

	s, err := NewState(ev, [4]int{25000, 25000, 25000, 25000})
	if err != nil {
		t.Fatalf("NewState() failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(*seat.MustSeat(0), *tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got := s.RenderBoard()
	want := "E-1 kyoku 0 honba  pipai: 69  dora_marker: E  \n" +
		"*{0} tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s 6m  \n" +
		"     ho:    \n" +
		" [1] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		" [2] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		" [3] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		"--------------------------------------------------------------------------------\n"
	if got != want {
		t.Errorf("RenderBoard() = %q, want %q", got, want)
	}
}

func TestState_RenderBoard_ActorAfterDiscard(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	discardedTile := *tile.MustTileFromCode("6m")

	if err := s.Apply(event.NewDraw(actor, discardedTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(actor, discardedTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}

	got := s.RenderBoard()
	want := "E-1 kyoku 0 honba  pipai: 69  dora_marker: E  \n" +
		"*{0} tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    6m \n" +
		" [1] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		" [2] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		" [3] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		"--------------------------------------------------------------------------------\n"
	if got != want {
		t.Errorf("RenderBoard() = %q, want %q", got, want)
	}
}

func TestState_RenderBoard_ActorAfterPon(t *testing.T) {
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

	got := s.RenderBoard()
	want := "E-1 kyoku 0 honba  pipai: 69  dora_marker: E  \n" +
		" {0} tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		" [1] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		" [2] tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s  \n" +
		"     ho:    \n" +
		"*[3] tehai: 1m 2m 3m 4m 5m 2p 3p 4p 2s 3s 4s  [1s(0)/1s 1s]\n" +
		"     ho:    \n" +
		"--------------------------------------------------------------------------------\n"
	if got != want {
		t.Errorf("RenderBoard() = %q, want %q", got, want)
	}
}
