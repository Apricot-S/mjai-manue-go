package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
)

type TenpaiProbEstimator struct {
	stats *configs.GameStats
}

func NewTenpaiProbEstimator(stats *configs.GameStats) *TenpaiProbEstimator {
	return &TenpaiProbEstimator{
		stats: stats,
	}
}

func (e *TenpaiProbEstimator) Estimate(player *game.Player, state game.State) float64 {
	if player.ReachState() != game.None {
		// If the player is in the riichi state, it is certain that the player is in tenpai.
		return 1.0
	}

	numRemainTurns := state.NumPipais() / 4
	numFuros := len(player.Furos())
	key := fmt.Sprintf("%d,%d", numRemainTurns, numFuros)
	if stat, found := e.stats.YamitenStats[key]; found {
		return float64(stat.Tenpai) / float64(stat.Total)
	}

	// If there is no stats, 1.0 is returned.
	return 1.0
}
