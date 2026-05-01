package round

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
)

func TestState_Apply_DrawRound(t *testing.T) {
	tests := []struct {
		name       string
		deltas     *[common.NumPlayers]int
		scores     *[common.NumPlayers]int
		wantScores [common.NumPlayers]int
	}{
		{
			name:       "scores",
			scores:     &[common.NumPlayers]int{23500, 26500, 23500, 26500},
			wantScores: [common.NumPlayers]int{23500, 26500, 23500, 26500},
		},
		{
			name:       "deltas",
			deltas:     &[common.NumPlayers]int{-1500, 1500, -1500, 1500},
			wantScores: [common.NumPlayers]int{23500, 26500, 23500, 26500},
		},
		{
			name:       "scores take precedence over deltas",
			deltas:     &[common.NumPlayers]int{1, 2, 3, 4},
			scores:     &[common.NumPlayers]int{23500, 26500, 23500, 26500},
			wantScores: [common.NumPlayers]int{23500, 26500, 23500, 26500},
		},
		{
			name:       "no scores or deltas",
			wantScores: [common.NumPlayers]int{25000, 25000, 25000, 25000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mustNewRoundStateForTest(t, newValidHands())

			if err := s.Apply(event.NewDrawRound(
				"fanpai",
				&[common.NumPlayers]bool{false, true, false, true},
				tt.deltas,
				tt.scores,
			)); err != nil {
				t.Fatalf("Apply(DrawRound) failed: %v", err)
			}

			if got := s.Scores(); got != tt.wantScores {
				t.Errorf("Scores() = %v, want %v", got, tt.wantScores)
			}
		})
	}
}
