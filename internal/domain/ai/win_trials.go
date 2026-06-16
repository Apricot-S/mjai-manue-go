package ai

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

const copiesPerTile = 4

func trialTileCounts(tiles []tile.Tile) hand.TileCounts34 {
	var counts hand.TileCounts34
	for _, t := range tiles {
		counts[t.RemoveRed().ID()]++
	}
	return counts
}

func wallTilesFromCounts(counts hand.TileCounts34) ([]tile.Tile, error) {
	for id, count := range counts {
		if count < 0 {
			return nil, fmt.Errorf("cannot build wall tiles: tile %s count must be non-negative", tile.MustTileFromID(id))
		}
	}
	return (&counts).ToTiles(), nil
}

func unseenWallFromVisibleTiles(visibleTiles []tile.Tile) ([]tile.Tile, error) {
	var counts hand.TileCounts34
	for id := range counts {
		counts[id] = copiesPerTile
	}
	for _, visible := range visibleTiles {
		if visible.IsUnknown() {
			return nil, fmt.Errorf("cannot build unseen wall: visible tile must not be unknown")
		}
		id := visible.RemoveRed().ID()
		counts[id]--
		if counts[id] < 0 {
			return nil, fmt.Errorf("cannot build unseen wall: tile %s is visible more than %d times", tile.MustTileFromID(id), copiesPerTile)
		}
	}
	return wallTilesFromCounts(counts)
}

func shuffledTrialTileCounts(wall []tile.Tile, numDraws int, rng *rand.Rand) (hand.TileCounts34, error) {
	if numDraws < 0 {
		return hand.TileCounts34{}, fmt.Errorf("cannot build trial tiles: numDraws must be non-negative")
	}
	if numDraws > len(wall) {
		return hand.TileCounts34{}, fmt.Errorf("cannot build trial tiles: numDraws %d exceeds wall length %d", numDraws, len(wall))
	}

	shuffled := slices.Clone(wall)
	rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return trialTileCounts(shuffled[:numDraws]), nil
}

func canAchieveGoalWithTrialTiles(goal service.Goal, trialTiles hand.TileCounts34) bool {
	for id, required := range goal.RequiredVector {
		if required > trialTiles[id] {
			return false
		}
	}
	return true
}

func trialWinPts(goals []winEstimateGoal, trialTiles hand.TileCounts34) (float64, bool, error) {
	best := 0.0
	for _, goal := range goals {
		if goal.points <= 0 {
			return 0, false, fmt.Errorf("cannot calculate trial win points: goal points must be positive")
		}
		if !canAchieveGoalWithTrialTiles(goal.Goal, trialTiles) {
			continue
		}
		best = max(best, goal.points)
	}
	return best, best > 0, nil
}

func candidateTrialWinPts(
	candidates []actionCandidate,
	goalsByKey map[string][]winEstimateGoal,
	trialTiles hand.TileCounts34,
) (map[string]float64, error) {
	if _, err := candidateTraceKeys(candidates); err != nil {
		return nil, err
	}

	winPtsByKey := make(map[string]float64)
	for _, candidate := range candidates {
		goals, ok := goalsByKey[candidate.traceKey]
		if !ok {
			return nil, fmt.Errorf("cannot calculate candidate trial win points: missing goals for %q", candidate.traceKey)
		}
		points, won, err := trialWinPts(goals, trialTiles)
		if err != nil {
			return nil, fmt.Errorf("cannot calculate candidate trial win points for %q: %w", candidate.traceKey, err)
		}
		if won {
			winPtsByKey[candidate.traceKey] = points
		}
	}
	return winPtsByKey, nil
}
