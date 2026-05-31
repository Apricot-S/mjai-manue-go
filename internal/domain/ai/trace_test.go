package ai

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func TestFormatDiscardTraceKey(t *testing.T) {
	discardTile := tile.MustTileFromCode("5m")
	if got := formatDiscardTraceKey(false, discardTile); got != "-1.5m" {
		t.Errorf("formatDiscardTraceKey(false) = %q, want %q", got, "-1.5m")
	}
	if got := formatDiscardTraceKey(true, discardTile); got != "0.5m" {
		t.Errorf("formatDiscardTraceKey(true) = %q, want %q", got, "0.5m")
	}
}

func TestFormatCandidateTrace(t *testing.T) {
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
			shanten:     1,
			score: candidateScore{
				averageRank:                 2.25,
				expectedPoints:              1200,
				dealInProb:                  0.125,
				winProb:                     0.25,
				exhaustiveDrawProb:          0.375,
				otherWinProb:                0.5,
				averageWinPoints:            3900,
				exhaustiveDrawAveragePoints: 1000,
			},
		},
	})
	want := "| action | avgRank | expPt | hojuProb | myHoraProb | ryukyokuProb | otherHoraProb | avgHoraPt | ryukyokuAvgPt | shanten | \n" +
		"|  -1.5m |  2.2500 |  1200 |    0.125 |      0.250 |        0.375 |         0.500 |      3900 |          1000 |       1 | \n"
	if got != want {
		t.Errorf("formatCandidateTrace() =\n%q\nwant\n%q", got, want)
	}
}

func TestFormatCandidateTrace_FormatsInfinityShanten(t *testing.T) {
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
			shanten:     service.InfinityShanten,
		},
	})
	if !strings.Contains(got, "Inf") {
		t.Errorf("formatCandidateTrace() = %q, want it to contain Inf", got)
	}
	if strings.Contains(got, fmt.Sprintf("%d", service.InfinityShanten)) {
		t.Errorf("formatCandidateTrace() = %q, should not contain raw InfinityShanten integer", got)
	}
}

func TestFormatTenpaiProbsTrace(t *testing.T) {
	got := formatTenpaiProbsTrace(
		[common.NumPlayers]float64{0.1, 0.25, 1, 0},
		seat.MustSeat(1),
	)
	want := "tenpaiProbs:  0: 0.100  2: 1.000  3: 0.000  \n"
	if got != want {
		t.Errorf("formatTenpaiProbsTrace() = %q, want %q", got, want)
	}
}

func TestFormatCandidateLog(t *testing.T) {
	self := seat.MustSeat(0)
	discard, err := action.NewDiscard(self, tile.MustTileFromCode("5m"), false)
	if err != nil {
		t.Fatalf("NewDiscard() failed: %v", err)
	}

	got := formatCandidateLog([]actionCandidate{
		{
			traceKey:    "-1.5m",
			action:      discard,
			discardTile: discard.Tile(),
			shanten:     1,
		},
	}, [common.NumPlayers]float64{0, 0.125, 0.5, 1}, self)
	if !strings.Contains(got, "\n\n\ntenpaiProbs:  1: 0.125  2: 0.500  3: 1.000  \n") {
		t.Errorf("formatCandidateLog() = %q, want tenpaiProbs after candidate table", got)
	}
}

func TestFormatDecisionTrace_AppendsDecidedKey(t *testing.T) {
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
		shanten:     1,
	}

	got := formatDecisionTrace(formatCandidateLog([]actionCandidate{*selected}, [common.NumPlayers]float64{}, self), selected)
	want := "| action | avgRank | expPt | hojuProb | myHoraProb | ryukyokuProb | otherHoraProb | avgHoraPt | ryukyokuAvgPt | shanten | \n" +
		"|  -1.5m |  0.0000 |     0 |    0.000 |      0.000 |        0.000 |         0.000 |         0 |             0 |       1 | \n" +
		"\n\n" +
		"tenpaiProbs:  1: 0.000  2: 0.000  3: 0.000  \n" +
		"decidedKey -1.5m\n"
	if got != want {
		t.Errorf("formatDecisionTrace() =\n%q\nwant\n%q", got, want)
	}
}
