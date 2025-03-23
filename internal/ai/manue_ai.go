package ai

import (
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

func NewManueAI(stats *configs.GameStats, root *configs.DangerNode) *ManueAI {
	return NewManueAIWithEstimators(
		stats,
		NewDangerEstimator(root),
		NewTenpaiProbEstimator(stats),
	)
}

func NewManueAIWithEstimators(stats *configs.GameStats, dangerEstimator *DangerEstimator, tenpaiProbEstimator *TenpaiProbEstimator) *ManueAI {
	return &ManueAI{
		stats:               stats,
		dangerEstimator:     dangerEstimator,
		tenpaiProbEstimator: tenpaiProbEstimator,
		noChanges:           [4]int{},
	}
}

func (a *ManueAI) DecideAction(state game.State, playerID int) (jsontext.Value, error) {
	panic("unimplemented!")
}
