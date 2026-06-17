package ai

import (
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestChooseBestCandidate_PrefersBlackTileOnTie(t *testing.T) {
	self := seat.MustSeat(0)
	redDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5mr"), false)
	if err != nil {
		t.Fatalf("NewDiscard(red) failed: %v", err)
	}
	blackDiscard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard(black) failed: %v", err)
	}

	candidates, err := buildSelfTurnCandidates([]action.Action{redDiscard, blackDiscard}, stubPlayerViewer{
		hand: hand.CodesToHand([]string{
			"1m", "2m", "3m", "4m", "5m", "5mr", "6m", "7m",
			"1p", "2p", "3p", "E", "E", "S",
		}),
		riichiState: player.NotRiichi,
		drawnTile:   nil,
	})
	if err != nil {
		t.Fatalf("buildSelfTurnCandidates() failed: %v", err)
	}
	evaluated := make([]evaluatedActionCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		evaluated = append(evaluated, evaluatedCandidateForTest(candidate, candidateScore{
			averageRank:    2.0,
			expectedPoints: 1000,
		}))
	}
	got := chooseBestCandidate(evaluated, true)
	if got.candidate.action != blackDiscard {
		t.Errorf("chooseBestCandidate() = %v, want black discard", got)
	}
}

func TestSortedCandidates_PrefersBlackForOrder(t *testing.T) {
	red := evaluatedCandidateForTest(actionCandidate{
		traceKey: "-1.5mr",
		red:      true,
	}, candidateScore{
		averageRank:    2.0,
		expectedPoints: 1000,
	})
	black := evaluatedCandidateForTest(actionCandidate{
		traceKey: "-1.5m",
		red:      false,
	}, candidateScore{
		averageRank:    2.0,
		expectedPoints: 1000,
	})

	got := sortedCandidates([]evaluatedActionCandidate{red, black}, true)
	if len(got) != 2 {
		t.Fatalf("len(sortedCandidates()) = %d, want 2", len(got))
	}
	if got[0].candidate.traceKey != black.candidate.traceKey {
		t.Errorf("first traceKey = %q, want black candidate %q", got[0].candidate.traceKey, black.candidate.traceKey)
	}
}

func TestSortedCandidates_CanIgnoreBlackPreference(t *testing.T) {
	red := evaluatedCandidateForTest(actionCandidate{
		traceKey: "-1.5mr",
		red:      true,
	}, candidateScore{
		averageRank:    2.0,
		expectedPoints: 1000,
	})
	black := evaluatedCandidateForTest(actionCandidate{
		traceKey: "-1.5m",
		red:      false,
	}, candidateScore{
		averageRank:    2.0,
		expectedPoints: 1000,
	})

	got := sortedCandidates([]evaluatedActionCandidate{red, black}, false)
	if len(got) != 2 {
		t.Fatalf("len(sortedCandidates()) = %d, want 2", len(got))
	}
	if got[0].candidate.traceKey != red.candidate.traceKey {
		t.Errorf("first traceKey = %q, want original first candidate %q", got[0].candidate.traceKey, red.candidate.traceKey)
	}
}

func TestCompareCandidates(t *testing.T) {
	tests := []struct {
		name        string
		lhs         evaluatedActionCandidate
		rhs         evaluatedActionCandidate
		preferBlack bool
		want        int
	}{
		{
			name:        "better average rank wins",
			lhs:         evaluatedCandidateForTest(actionCandidate{}, candidateScore{averageRank: 1.9, expectedPoints: 0}),
			rhs:         evaluatedCandidateForTest(actionCandidate{}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "higher expected points wins on rank tie",
			lhs:         evaluatedCandidateForTest(actionCandidate{}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			rhs:         evaluatedCandidateForTest(actionCandidate{}, candidateScore{averageRank: 2.0, expectedPoints: 900}),
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "black tile wins on expected value tie",
			lhs:         evaluatedCandidateForTest(actionCandidate{red: false}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			rhs:         evaluatedCandidateForTest(actionCandidate{red: true}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "red tie ignored when preferBlack is false",
			lhs:         evaluatedCandidateForTest(actionCandidate{red: false}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			rhs:         evaluatedCandidateForTest(actionCandidate{red: true}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			preferBlack: false,
			want:        0,
		},
		{
			name:        "complete tie returns zero",
			lhs:         evaluatedCandidateForTest(actionCandidate{}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			rhs:         evaluatedCandidateForTest(actionCandidate{}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
			preferBlack: false,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareCandidates(tt.lhs, tt.rhs, tt.preferBlack)
			if got != tt.want {
				t.Errorf("compareCandidates() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestChooseBestCandidate_DoesNotPreferRiichiOnScoreTie(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	candidates := []evaluatedActionCandidate{
		evaluatedCandidateForTest(actionCandidate{action: discard, riichi: false}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
		evaluatedCandidateForTest(actionCandidate{action: riichi, riichi: true}, candidateScore{averageRank: 2.0, expectedPoints: 1000}),
	}

	got := chooseBestCandidate(candidates, false)
	if got.candidate.action != discard {
		t.Errorf("chooseBestCandidate() = %T %[1]v, want first tied candidate", got.candidate.action)
	}
}

func evaluatedCandidateForTest(candidate actionCandidate, score candidateScore) evaluatedActionCandidate {
	return evaluatedActionCandidate{
		candidate: candidate,
		score:     score,
	}
}
