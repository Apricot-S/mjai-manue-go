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
	hc, err := state.HoraCandidate()
	if err != nil {
		return nil, err
	}
	if hc != nil {
		// If it can win, always win
		hora, err := message.NewHora(playerID, hc.Target(), hc.Pai().ToString(), 0, nil, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create hora message: %w", err)
		}
		res, err := json.Marshal(&hora)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal hora message: %w", err)
		}
		return res, nil
	}

	decision, err := a.decideDahai(state, playerID)
	if err != nil {
		return nil, err
	}
	if decision != nil {
		return decision, nil
	}

	decision, err = a.decideFuro(state, playerID)
	if err != nil {
		return nil, err
	}
	if decision != nil {
		return decision, nil
	}

	// no action is possible
	none := message.NewNone()
	res, err := json.Marshal(&none)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal none message: %w", err)
	}
	return res, nil
}

func (a *ManueAI) decideDahai(state game.StateAnalyzer, playerID int) (jsontext.Value, error) {
	dc := state.DahaiCandidates()
	rdc, err := state.ReachDahaiCandidates()
	fd := state.ForbiddenDahais()
	if err != nil {
		return nil, err
	}
	if len(dc) == 0 && len(rdc) == 0 {
		// no action is possible
		none := message.NewNone()
		res, err := json.Marshal(&none)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal none message: %w", err)
		}
		return res, nil
	}

	// my turn
	if state.Players()[playerID].ReachState() == game.Accepted {
		// in reach
		dahai, err := message.NewDahai(playerID, dc[0].ToString(), true, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create dahai message: %w", err)
		}
		res, err := json.Marshal(&dahai)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal dahai message: %w", err)
		}
		return res, nil
	}

	pai, isReach, err := a.getMetrics(dc, rdc, fd)
	if err != nil {
		return nil, err
	}

	if isReach {
		// reach declaration
		reach, err := message.NewReach(playerID, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create reach message: %w", err)
		}
		res, err := json.Marshal(&reach)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal reach message: %w", err)
		}
		return res, nil
	}

	// dahai
	isTsumogiri := state.IsTsumoPai(pai)
	dahai, err := message.NewDahai(playerID, pai.ToString(), isTsumogiri, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create dahai message: %w", err)
	}
	res, err := json.Marshal(&dahai)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dahai message: %w", err)
	}
	return res, nil
}

func (a *ManueAI) decideFuro(state game.StateAnalyzer, playerID int) (jsontext.Value, error) {
	fc, err := state.FuroCandidates()
	if err != nil {
		return nil, err
	}
	if len(fc) == 0 {
		// no action is possible
		none := message.NewNone()
		res, err := json.Marshal(&none)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal none message: %w", err)
		}
		return res, nil
	}

	panic("unimplemented")
}
