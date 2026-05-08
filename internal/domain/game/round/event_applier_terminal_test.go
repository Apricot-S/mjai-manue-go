package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_AfterWinAllowsOnlyAdditionalWins(t *testing.T) {
	s := newStateAfterRonForTerminalTest(t, seat.MustSeat(1))
	target := seat.MustSeat(0)
	winningTile := tile.MustTileFromCode("6m")

	if err := s.Apply(event.NewDraw(seat.MustSeat(2), tile.MustTileFromCode("7m"))); err == nil {
		t.Fatal("Apply(Draw) after Win succeeded unexpectedly")
	}
	if err := s.Apply(event.NewWin(
		seat.MustSeat(2),
		target,
		&winningTile,
		8000,
		nil,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
	)); err != nil {
		t.Fatalf("Apply(additional Win) failed: %v", err)
	}
}

func TestState_Apply_AfterDrawRoundReturnsError(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	if err := s.Apply(event.NewDrawRound("fanpai", nil, nil, nil)); err != nil {
		t.Fatalf("Apply(DrawRound) failed: %v", err)
	}

	if err := s.Apply(event.NewDraw(seat.MustSeat(0), tile.MustTileFromCode("6m"))); err == nil {
		t.Fatal("Apply(Draw) after DrawRound succeeded unexpectedly")
	}
	if err := s.Apply(event.NewWin(
		seat.MustSeat(1),
		seat.MustSeat(0),
		new(tile.MustTileFromCode("6m")),
		8000,
		nil,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
	)); err == nil {
		t.Fatal("Apply(Win) after DrawRound succeeded unexpectedly")
	}
}

func TestState_LegalActions_AfterTerminalEventIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		s    *State
	}{
		{
			name: "win",
			s:    newStateAfterRonForTerminalTest(t, seat.MustSeat(1)),
		},
		{
			name: "draw round",
			s: func() *State {
				s := mustNewRoundStateForTest(t, newValidHands())
				if err := s.Apply(event.NewDrawRound("fanpai", nil, nil, nil)); err != nil {
					t.Fatalf("Apply(DrawRound) failed: %v", err)
				}
				return s
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.LegalActions(seat.MustSeat(0))
			if err != nil {
				t.Fatalf("LegalActions() failed: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("LegalActions() = %v, want empty", got)
			}
		})
	}
}

func newStateAfterRonForTerminalTest(t *testing.T, winActor seat.Seat) *State {
	t.Helper()

	hands := newValidHands()
	hands[winActor.Index()] = tenpaiHandWaiting36mForTest()
	// The terminal tests also apply an additional ron by player 2.
	hands[seat.MustSeat(2).Index()] = tenpaiHandWaiting36mForTest()
	s := mustNewRoundStateForTest(t, hands)
	target := seat.MustSeat(0)
	winningTile := tile.MustTileFromCode("6m")
	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		winActor,
		target,
		&winningTile,
		8000,
		nil,
		&[common.NumPlayers]int{25000, 25000, 25000, 25000},
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}
	return s
}
