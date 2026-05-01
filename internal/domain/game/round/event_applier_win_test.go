package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestState_Apply_Win(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	scores := [common.NumPlayers]int{25000, 30800, 34700, 9500}

	if err := s.Apply(event.NewWin(
		*seat.MustSeat(2),
		*seat.MustSeat(3),
		tile.MustTileFromCode("9m"),
		8000,
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(Win) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}
