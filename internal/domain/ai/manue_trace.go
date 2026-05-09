package ai

import (
	"fmt"
	"slices"
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

func formatDecisionTrace(candidates []actionCandidate, selected *actionCandidate) string {
	trace := formatCandidateLog(candidates)
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
	if len(candidates) == 0 {
		return ""
	}

	sortedCandidates := slices.Clone(candidates)
	slices.SortFunc(sortedCandidates, func(lhs, rhs actionCandidate) int {
		return compareCandidateScore(&lhs.score, &rhs.score, true)
	})

	rows := [][]string{{
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
	}}
	for _, candidate := range sortedCandidates {
		rows = append(rows, []string{
			candidate.traceKey,
			fmt.Sprintf("%.4f", candidate.score.avgRank),
			fmt.Sprintf("%.0f", candidate.score.expPts),
			fmt.Sprintf("%.3f", candidate.score.dealInProb),
			fmt.Sprintf("%.3f", candidate.score.winProb),
			fmt.Sprintf("%.3f", candidate.score.drawProb),
			fmt.Sprintf("%.3f", candidate.score.othersWinProb),
			fmt.Sprintf("%.0f", candidate.score.avgWinPts),
			fmt.Sprintf("%.0f", candidate.score.avgDrawPts),
			formatShantenTraceValue(candidate.score.shanten),
		})
	}
	return formatTraceTable(rows)
}

func formatShantenTraceValue(shanten int) string {
	if shanten == service.InfinityShanten {
		return "Infinity"
	}
	return fmt.Sprintf("%d", shanten)
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
