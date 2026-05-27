package ai

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

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
