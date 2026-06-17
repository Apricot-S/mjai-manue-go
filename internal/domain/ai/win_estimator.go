package ai

import (
	"fmt"
	"math/rand/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
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
