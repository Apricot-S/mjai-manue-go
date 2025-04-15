package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/go-json-experiment/json/jsontext"
)

type ManueAI struct {
	stats               *configs.GameStats
	dangerEstimator     *DangerEstimator
	tenpaiProbEstimator *TenpaiProbEstimator
	noChanges           [4]int
}

func NewManueAI() (*ManueAI, error) {
	stats, err := configs.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get the stats: %w", err)
	}
	root, err := configs.GetDangerTree()
	if err != nil {
		return nil, fmt.Errorf("failed to get the danger tree: %w", err)
	}

	return NewManueAIWithEstimators(
		stats,
		NewDangerEstimator(root),
		NewTenpaiProbEstimator(stats),
	), nil
}

func NewManueAIWithEstimators(
	stats *configs.GameStats,
	dangerEstimator *DangerEstimator,
	tenpaiProbEstimator *TenpaiProbEstimator,
) *ManueAI {
	return &ManueAI{
		stats:               stats,
		dangerEstimator:     dangerEstimator,
		tenpaiProbEstimator: tenpaiProbEstimator,
		noChanges:           [4]int{},
	}
}

func (a *ManueAI) DecideAction(state game.StateAnalyzer, playerID int) (jsontext.Value, error) {
	panic("unimplemented!")
}
