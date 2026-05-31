package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service/block"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type winEstimateGoal struct {
	service.Goal
	points float64
}

const (
	numRedFiveTileTypes = 3
	winningHandSize     = 14
	shantenPruneLimit   = 3
)

func filteredWinEstimateGoals(candidate actionCandidate) []service.Goal {
	goals := make([]service.Goal, 0, len(candidate.shantenGoals))
	for _, goal := range candidate.shantenGoals {
		if candidate.scoreAsRiichi && goal.Shanten > 0 {
			continue
		}
		if candidate.baseShanten > shantenPruneLimit && goal.Shanten > candidate.baseShanten {
			continue
		}
		if !candidate.discardTile.IsUnknown() {
			discardID := candidate.discardTile.RemoveRed().ID()
			if goal.ThrowableVector[discardID] <= 0 {
				continue
			}
		}
		goals = append(goals, goal)
	}
	return goals
}

func scoredWinEstimateGoals(candidate actionCandidate, context winEstimateGoalContext) ([]winEstimateGoal, error) {
	if candidate.turnHand == nil {
		return nil, fmt.Errorf("cannot score win estimate goals: turn hand must not be nil")
	}

	goals := filteredWinEstimateGoals(candidate)
	scoredGoals := make([]winEstimateGoal, 0, len(goals))
	for _, goal := range goals {
		scoringHand, err := scoringHandForGoal(candidate.turnHand, goal.Blocks)
		if err != nil {
			return nil, err
		}
		fu, han, _ := service.CalculateFuHan(
			scoringHand,
			goal.Blocks,
			context.melds,
			context.roundWind,
			context.seatWind,
			context.doraIndicators,
			candidate.scoreAsRiichi,
		)
		points := service.RonPoints(fu, han, context.dealer)
		if points <= 0 {
			continue
		}
		scoredGoals = append(scoredGoals, winEstimateGoal{
			Goal:   goal,
			points: float64(points),
		})
	}
	return scoredGoals, nil
}

func scoringHandForGoal(sourceHand *hand.VisibleHand, blocks []block.Block) (*hand.VisibleHand, error) {
	redCounts := make(map[int]int, numRedFiveTileTypes)
	for _, t := range sourceHand.ToTiles() {
		if t.IsRed() {
			redCounts[t.RemoveRed().ID()]++
		}
	}

	tiles := make([]tile.Tile, 0, winningHandSize)
	for _, b := range blocks {
		for _, t := range b.ToTiles() {
			normal := t.RemoveRed()
			if redCounts[normal.ID()] > 0 {
				tiles = append(tiles, normal.AddRed())
				redCounts[normal.ID()]--
				continue
			}
			tiles = append(tiles, normal)
		}
	}

	scoringHand, err := hand.NewVisibleHand(tiles)
	if err != nil {
		return nil, fmt.Errorf("cannot build scoring hand for win estimate goal: %w", err)
	}
	return scoringHand, nil
}

func scoredWinEstimateGoalsByKey(
	candidates []actionCandidate,
	context winEstimateGoalContext,
) (map[string][]winEstimateGoal, error) {
	if _, err := candidateTraceKeys(candidates); err != nil {
		return nil, err
	}

	goalsByKey := make(map[string][]winEstimateGoal, len(candidates))
	for _, candidate := range candidates {
		candidateContext := context
		if candidate.melds != nil {
			candidateContext.melds = candidate.melds
		}
		goals, err := scoredWinEstimateGoals(candidate, candidateContext)
		if err != nil {
			return nil, fmt.Errorf("cannot score win estimate goals for %q: %w", candidate.traceKey, err)
		}
		goalsByKey[candidate.traceKey] = goals
	}
	return goalsByKey, nil
}

func countWinEstimateGoalsByGroup(candidates []actionCandidate, goalsByKey map[string][]winEstimateGoal) []int {
	maxGroup := 0
	for _, candidate := range candidates {
		maxGroup = max(maxGroup, candidate.evaluationGroup)
	}
	counts := make([]int, maxGroup+1)
	for _, candidate := range candidates {
		counts[candidate.evaluationGroup] += len(goalsByKey[candidate.traceKey])
	}
	return counts
}
