package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

func TestState_Apply_DrawRound(t *testing.T) {
	s := mustNewRoundStateForTest(t, newValidHands())
	scores := [common.NumPlayers]int{23500, 26500, 23500, 26500}

	if err := s.Apply(event.NewDrawRound(
		"fanpai",
		&[common.NumPlayers]bool{false, true, false, true},
		nil,
		&scores,
	)); err != nil {
		t.Fatalf("Apply(DrawRound) failed: %v", err)
	}

	if got := s.Scores(); got != scores {
		t.Errorf("Scores() = %v, want %v", got, scores)
	}
}
