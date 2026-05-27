package ai

import (
	"fmt"
	"math/rand/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type winEstimateGoalContext struct {
	melds          []meld.Meld
	roundWind      wind.Wind
	seatWind       wind.Wind
	doraIndicators []tile.Tile
	dealer         bool
}

type winEstimateStateViewer interface {
	VisibleTiles(playerSeat seat.Seat) tile.Tiles
	Turn() float64
}

func candidateTraceKeys(candidates []actionCandidate) ([]string, error) {
	keys := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate.traceKey == "" {
			return nil, fmt.Errorf("cannot build candidate keys: trace key must not be empty")
		}
		if _, ok := seen[candidate.traceKey]; ok {
			return nil, fmt.Errorf("cannot build candidate keys: duplicate trace key %q", candidate.traceKey)
		}
		seen[candidate.traceKey] = struct{}{}
		keys = append(keys, candidate.traceKey)
	}
	return keys, nil
}

func winEstimatesForCandidates(candidates []actionCandidate, trials []map[string]float64) (map[string]winEstimate, error) {
	keys, err := candidateTraceKeys(candidates)
	if err != nil {
		return nil, err
	}
	return winEstimatesFromTrials(keys, trials)
}

func winEstimatesFromTrialTiles(
	candidates []actionCandidate,
	goalsByKey map[string][]winEstimateGoal,
	trials [][]tile.Tile,
) (map[string]winEstimate, error) {
	trialResults := make([]map[string]float64, 0, len(trials))
	for i, trial := range trials {
		points, err := candidateTrialWinPts(candidates, goalsByKey, trialTileCounts(trial))
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimates from trial %d: %w", i, err)
		}
		trialResults = append(trialResults, points)
	}
	return winEstimatesForCandidates(candidates, trialResults)
}

func winEstimatesFromShuffledWall(
	candidates []actionCandidate,
	goalsByKey map[string][]winEstimateGoal,
	wall []tile.Tile,
	numDraws int,
	numTries int,
	rng *rand.Rand,
) (map[string]winEstimate, error) {
	keys, err := candidateTraceKeys(candidates)
	if err != nil {
		return nil, err
	}
	if numTries <= 0 {
		return nil, fmt.Errorf("cannot build win estimates from shuffled wall: numTries must be positive")
	}
	if rng == nil {
		return nil, fmt.Errorf("cannot build win estimates from shuffled wall: rng must not be nil")
	}

	accumulators := newWinEstimateAccumulatorSet(keys)
	for i := range numTries {
		trialTiles, err := shuffledTrialTileCounts(wall, numDraws, rng)
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimates from shuffled wall trial %d: %w", i, err)
		}
		points, err := candidateTrialWinPts(candidates, goalsByKey, trialTiles)
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimates from shuffled wall trial %d: %w", i, err)
		}
		if err := accumulators.addTrial(points); err != nil {
			return nil, fmt.Errorf("cannot build win estimates from shuffled wall trial %d: %w", i, err)
		}
	}
	return accumulators.estimates()
}

func winEstimatesFromState(
	stats RoundEndStats,
	state winEstimateStateViewer,
	self seat.Seat,
	candidates []actionCandidate,
	goalsByKey map[string][]winEstimateGoal,
	numTries int,
	rng *rand.Rand,
) (map[string]winEstimate, error) {
	numDraws, err := expectedRemainingTurns(stats, state.Turn())
	if err != nil {
		return nil, err
	}
	wall, err := unseenWallFromVisibleTiles(state.VisibleTiles(self))
	if err != nil {
		return nil, err
	}
	return winEstimatesFromShuffledWall(candidates, goalsByKey, wall, numDraws, numTries, rng)
}

func candidateShanten(discardTile tile.Tile, baseShanten int, goals []service.Goal) int {
	if discardTile.IsUnknown() {
		return baseShanten
	}
	discardID := discardTile.RemoveRed().ID()
	shanten := service.InfinityShanten
	for _, goal := range goals {
		if goal.ThrowableVector[discardID] > 0 && goal.Shanten < shanten {
			shanten = goal.Shanten
		}
	}
	return shanten
}
