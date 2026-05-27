package ai

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestManueAgent_chooseBestCandidate_PrefersBlackTileOnTie(t *testing.T) {
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
	got := chooseBestCandidate(candidates, true)
	if got.action != blackDiscard {
		t.Errorf("chooseBestCandidate() = %v, want black discard", got)
	}
}

func TestManueAgent_formatDiscardTraceKey(t *testing.T) {
	discardTile := tile.MustTileFromCode("5m")
	if got := formatDiscardTraceKey(false, discardTile); got != "-1.5m" {
		t.Errorf("formatDiscardTraceKey(false) = %q, want %q", got, "-1.5m")
	}
	if got := formatDiscardTraceKey(true, discardTile); got != "0.5m" {
		t.Errorf("formatDiscardTraceKey(true) = %q, want %q", got, "0.5m")
	}
}

func TestManueAgent_formatCandidateTrace(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got := formatCandidateTrace([]actionCandidate{
		{
			traceKey:    "-1.5m",
			action:      discard,
			riichi:      false,
			discardTile: discard.Tile(),
			score: candidateScore{
				averageRank:                 2.25,
				expectedPoints:              1200,
				dealInProb:                  0.125,
				winProb:                     0.25,
				exhaustiveDrawProb:          0.375,
				otherWinProb:                0.5,
				averageWinPoints:            3900,
				exhaustiveDrawAveragePoints: 1000,
				shanten:                     1,
			},
		},
	})
	want := "| action | avgRank | expPt | hojuProb | myHoraProb | ryukyokuProb | otherHoraProb | avgHoraPt | ryukyokuAvgPt | shanten | \n" +
		"|  -1.5m |  2.2500 |  1200 |    0.125 |      0.250 |        0.375 |         0.500 |      3900 |          1000 |       1 | \n"
	if got != want {
		t.Errorf("formatCandidateTrace() =\n%q\nwant\n%q", got, want)
	}
}

func TestManueAgent_formatCandidateTrace_FormatsInfinityShanten(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got := formatCandidateTrace([]actionCandidate{
		{
			traceKey:    "-1.5m",
			action:      discard,
			riichi:      false,
			discardTile: discard.Tile(),
			score: candidateScore{
				shanten: service.InfinityShanten,
			},
		},
	})
	if !strings.Contains(got, "Inf") {
		t.Errorf("formatCandidateTrace() = %q, want it to contain Inf", got)
	}
	if strings.Contains(got, fmt.Sprintf("%d", service.InfinityShanten)) {
		t.Errorf("formatCandidateTrace() = %q, should not contain raw InfinityShanten integer", got)
	}
}

func TestManueAgent_formatDecisionTrace_AppendsDecidedKey(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}
	selected := &actionCandidate{
		traceKey:    "-1.5m",
		action:      discard,
		riichi:      false,
		discardTile: discard.Tile(),
		score: candidateScore{
			shanten: 1,
		},
	}

	got := formatDecisionTrace(formatCandidateLog([]actionCandidate{*selected}), selected)
	if !strings.HasSuffix(got, "\n\n\ndecidedKey -1.5m\n") {
		t.Errorf("formatDecisionTrace() = %q, want two blank lines before decidedKey suffix", got)
	}
}

func TestSortedTraceCandidates_PrefersBlackForDisplayOrder(t *testing.T) {
	red := actionCandidate{
		traceKey: "-1.5mr",
		score: candidateScore{
			averageRank:    2.0,
			expectedPoints: 1000,
			red:            true,
		},
	}
	black := actionCandidate{
		traceKey: "-1.5m",
		score: candidateScore{
			averageRank:    2.0,
			expectedPoints: 1000,
			red:            false,
		},
	}

	got := sortedTraceCandidates([]actionCandidate{red, black})
	if len(got) != 2 {
		t.Fatalf("len(sortedTraceCandidates()) = %d, want 2", len(got))
	}
	if got[0].traceKey != black.traceKey {
		t.Errorf("first traceKey = %q, want black candidate %q", got[0].traceKey, black.traceKey)
	}
}

func TestManueAgent_compareCandidateScore(t *testing.T) {
	tests := []struct {
		name        string
		lhs         candidateScore
		rhs         candidateScore
		preferBlack bool
		want        int
	}{
		{
			name:        "better average rank wins",
			lhs:         candidateScore{averageRank: 1.9, expectedPoints: 0},
			rhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000},
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "higher expected points wins on rank tie",
			lhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000},
			rhs:         candidateScore{averageRank: 2.0, expectedPoints: 900},
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "black tile wins on expected value tie",
			lhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000, red: false},
			rhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000, red: true},
			preferBlack: true,
			want:        -1,
		},
		{
			name:        "red tie ignored when preferBlack is false",
			lhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000, red: false},
			rhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000, red: true},
			preferBlack: false,
			want:        0,
		},
		{
			name:        "complete tie returns zero",
			lhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000},
			rhs:         candidateScore{averageRank: 2.0, expectedPoints: 1000},
			preferBlack: false,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareCandidateScore(&tt.lhs, &tt.rhs, tt.preferBlack)
			if got != tt.want {
				t.Errorf("compareCandidateScore() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestManueAgent_chooseBestCandidate_DoesNotPreferRiichiOnScoreTie(t *testing.T) {
	self := seat.MustSeat(0)
	riichi := action.NewRiichi(self)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	candidates := []actionCandidate{
		{action: discard, riichi: false, score: candidateScore{averageRank: 2.0, expectedPoints: 1000}},
		{action: riichi, riichi: true, score: candidateScore{averageRank: 2.0, expectedPoints: 1000}},
	}

	got := chooseBestCandidate(candidates, false)
	if got.action != discard {
		t.Errorf("chooseBestCandidate() = %T %[1]v, want first tied candidate", got.action)
	}
}
