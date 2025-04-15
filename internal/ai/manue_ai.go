package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
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
	dc, err := state.DahaiCandidates()
	if err != nil {
		return nil, err
	}
	rdc, err := state.ReachDahaiCandidates()
	if err != nil {
		return nil, err
	}
	cc, err := state.ChiCandidates()
	if err != nil {
		return nil, err
	}
	pc, err := state.PonCandidates()
	if err != nil {
		return nil, err
	}
	canHora, err := state.CanHora()
	if err != nil {
		return nil, err
	}

	isMyTurn := dc != nil || rdc != nil || canHora
	canCallOrRon := cc != nil || pc != nil || canHora

	if !isMyTurn && !canCallOrRon {
		// no action is possible
		none := message.NewNone()
		res, err := json.Marshal(&none)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal none message: %w", err)
		}
		return res, nil
	}

	panic("unimplemented!")
}
