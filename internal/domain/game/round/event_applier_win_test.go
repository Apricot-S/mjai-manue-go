package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Win(t *testing.T) {
	tests := []struct {
		name       string
		deltas     *[common.NumPlayers]int
		scores     *[common.NumPlayers]int
		wantScores [common.NumPlayers]int
	}{
		{
			name:       "scores",
			scores:     &[common.NumPlayers]int{73000, 9000, 9000, 9000},
			wantScores: [common.NumPlayers]int{73000, 9000, 9000, 9000},
		},
		{
			name:       "deltas",
			deltas:     &[common.NumPlayers]int{48000, -16000, -16000, -16000},
			wantScores: [common.NumPlayers]int{73000, 9000, 9000, 9000},
		},
		{
			name:       "scores take precedence over deltas",
			deltas:     &[common.NumPlayers]int{1, 2, 3, 4},
			scores:     &[common.NumPlayers]int{73000, 9000, 9000, 9000},
			wantScores: [common.NumPlayers]int{73000, 9000, 9000, 9000},
		},
		{
			name:       "no scores or deltas",
			wantScores: [common.NumPlayers]int{25000, 25000, 25000, 25000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mustNewRoundStateForTest(t, newValidHands())
			actor := *seat.MustSeat(0)
			winningTile := *tile.MustTileFromCode("6m")

			if err := s.Apply(event.NewDraw(actor, winningTile)); err != nil {
				t.Fatalf("Apply(Draw) failed: %v", err)
			}
			if err := s.Apply(event.NewWin(
				actor,
				actor,
				&winningTile,
				48000,
				tt.deltas,
				tt.scores,
			)); err != nil {
				t.Fatalf("Apply(Win) failed: %v", err)
			}

			if got := s.Scores(); got != tt.wantScores {
				t.Errorf("Scores() = %v, want %v", got, tt.wantScores)
			}
		})
	}
}

func TestState_Apply_Win_ReturnsErrorBeforeFirstDraw(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	scores := [common.NumPlayers]int{25000, 30800, 34700, 9500}

	if err := s.Apply(event.NewWin(
		*seat.MustSeat(2),
		*seat.MustSeat(3),
		tile.MustTileFromCode("9m"),
		8000,
		nil,
		&scores,
	)); err == nil {
		t.Fatal("Apply(Win) succeeded unexpectedly")
	}

	if got := s.Scores(); got != [common.NumPlayers]int{25000, 25000, 25000, 25000} {
		t.Errorf("Scores() = %v, want unchanged initial scores", got)
	}
}

func TestState_Apply_Win_Renhou(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(1)
	target := *seat.MustSeat(0)
	winningTile := *tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{57000, -7000, 25000, 25000}

	if err := s.Apply(event.NewDraw(target, winningTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewDiscard(target, winningTile, true)); err != nil {
		t.Fatalf("Apply(Discard) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		target,
		&winningTile,
		32000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}

func TestState_Apply_Win_TsumoWithoutWinningTile(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	actor := *seat.MustSeat(0)
	drawnTile := *tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{73000, 9000, 9000, 9000}

	if err := s.Apply(event.NewDraw(actor, drawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		actor,
		nil,
		48000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}

func TestState_Apply_Win_InvisibleTsumo(t *testing.T) {
	hands := [common.NumPlayers][common.InitHandSize]tile.Tile{}
	for p := range common.NumPlayers {
		for i := range common.InitHandSize {
			hands[p][i] = *tile.MustTileFromCode("?")
		}
	}
	s := mustNewRoundStateForTest(t, hands)
	actor := *seat.MustSeat(2)
	unknownDrawnTile := *tile.MustTileFromCode("?")
	winningTile := *tile.MustTileFromCode("6m")
	scores := [common.NumPlayers]int{9000, 9000, 73000, 9000}

	if err := s.Apply(event.NewDraw(actor, unknownDrawnTile)); err != nil {
		t.Fatalf("Apply(Draw) failed: %v", err)
	}
	if err := s.Apply(event.NewWin(
		actor,
		actor,
		&winningTile,
		48000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}
