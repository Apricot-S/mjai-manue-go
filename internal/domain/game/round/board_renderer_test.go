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
	ev, err := event.NewStartRound(
		wind.East,
		1,
		0,
		0,
		*seat.MustSeat(0),
		*tile.MustTileFromCode("E"),
		&[4]int{25000, 25000, 25000, 25000},
		hands,
	)
	if err != nil {
		t.Fatalf("event.NewStartRound() failed: %v", err)
	}

	s, err := NewState(ev, [4]int{25000, 25000, 25000, 25000})
	if err != nil {
		t.Fatalf("NewState() failed: %v", err)
	}
	if err := s.Apply(event.NewDraw(*seat.MustSeat(0), *tile.MustTileFromCode("6m"))); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}

	got := s.RenderBoard()
	want := "E-1 kyoku 0 honba  pipai: 69  dora_marker: E  \n" +
		" {0} tehai: 1m 2m 3m 4m 5m 1p 2p 3p 4p 1s 2s 3s 4s 6m  \n" +
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
