package ai

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func formatDiscardTraceKey(riichi bool, discardTile tile.Tile) string {
	prefix := -1
	if riichi {
		prefix = 0
	}
	return fmt.Sprintf("%d.%s", prefix, discardTile)
}

func formatDecisionTrace(log string, selected *actionCandidate, summary candidateEvaluationSummary) string {
	var b strings.Builder
	b.WriteString(formatGoalCountTrace(summary))
	b.WriteString(log)
	if selected == nil {
		return b.String()
	}
	fmt.Fprintf(&b, "decidedKey %s\n", selected.traceKey)
	return b.String()
}

func formatGoalCountTrace(summary candidateEvaluationSummary) string {
	var b strings.Builder
	for _, count := range summary.winEstimateGoalCounts {
		fmt.Fprintf(&b, "goals %d\n", count)
	}
	return b.String()
}

func formatCandidateLog(candidates []actionCandidate, tenpaiProbs [common.NumPlayers]float64, self seat.Seat) string {
	trace := formatCandidateTrace(candidates)
	if trace == "" {
		return ""
	}
	return trace + "\n\n" + formatTenpaiProbsTrace(tenpaiProbs, self)
}

func formatTenpaiProbsTrace(tenpaiProbs [common.NumPlayers]float64, self seat.Seat) string {
	var b strings.Builder
	b.WriteString("tenpaiProbs:  ")
	for i := range common.NumPlayers {
		if seat.MustSeat(i) == self {
			continue
		}
		fmt.Fprintf(&b, "%d: %.3f  ", i, tenpaiProbs[i])
	}
	b.WriteString("\n")
	return b.String()
}

func formatCandidateTrace(candidates []actionCandidate) string {
	n := len(candidates)
	if n == 0 {
		return ""
	}

	rows := make([][]string, n+1)
	rows[0] = []string{
		"action",
		"avgRank",
		"expPt",
		"hojuProb",
		"myHoraProb",
		"ryukyokuProb",
		"otherHoraProb",
		"avgHoraPt",
		"ryukyokuAvgPt",
		"shanten",
	}
	for i, candidate := range sortedCandidates(candidates, true) {
		rows[i+1] = []string{
			candidate.traceKey,
			strconv.FormatFloat(candidate.score.averageRank, 'f', 4, 64),
			strconv.FormatFloat(candidate.score.expectedPoints, 'f', 0, 64),
			strconv.FormatFloat(candidate.score.dealInProb, 'f', 3, 64),
			strconv.FormatFloat(candidate.score.winProb, 'f', 3, 64),
			strconv.FormatFloat(candidate.score.exhaustiveDrawProb, 'f', 3, 64),
			strconv.FormatFloat(candidate.score.otherWinProb, 'f', 3, 64),
			strconv.FormatFloat(candidate.score.averageWinPoints, 'f', 0, 64),
			strconv.FormatFloat(candidate.score.exhaustiveDrawAveragePoints, 'f', 0, 64),
			formatShantenTraceValue(candidate.shanten),
		}
	}
	return formatTraceTable(rows)
}

func formatShantenTraceValue(shanten int) string {
	if shanten == service.InfinityShanten {
		return "Inf"
	}
	return strconv.Itoa(shanten)
}

func formatTraceTable(rows [][]string) string {
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			widths[i] = max(widths[i], len(cell))
		}
	}

	var b strings.Builder
	for _, row := range rows {
		b.WriteString("| ")
		for i, cell := range row {
			fmt.Fprintf(&b, "%*s | ", widths[i], cell)
		}
		b.WriteString("\n")
	}
	return b.String()
}
