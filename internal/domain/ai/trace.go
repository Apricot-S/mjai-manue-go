package ai

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func formatDiscardTraceKey(riichi bool, discardTile tile.Tile) string {
	prefix := -1
	if riichi {
		prefix = 0
	}
	return fmt.Sprintf("%d.%s", prefix, discardTile)
}

func formatDecisionTrace(log string, selected *actionCandidate) string {
	trace := log
	if selected == nil {
		return trace
	}
	return trace + fmt.Sprintf("decidedKey %s\n", selected.traceKey)
}

func formatCandidateLog(candidates []actionCandidate) string {
	trace := formatCandidateTrace(candidates)
	if trace == "" {
		return ""
	}
	return trace + "\n\n"
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
	for i, candidate := range sortedTraceCandidates(candidates) {
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
			formatShantenTraceValue(candidate.score.shanten),
		}
	}
	return formatTraceTable(rows)
}

func sortedTraceCandidates(candidates []actionCandidate) []actionCandidate {
	sortedCandidates := slices.Clone(candidates)
	slices.SortFunc(sortedCandidates, func(lhs, rhs actionCandidate) int {
		return compareCandidateScore(&lhs.score, &rhs.score, true)
	})
	return sortedCandidates
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
