package ai

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/Apricot-S/mjai-manue-go/configs"
	"github.com/Apricot-S/mjai-manue-go/internal/ai/estimator"
	"github.com/Apricot-S/mjai-manue-go/internal/game"
	"github.com/Apricot-S/mjai-manue-go/internal/message"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type ManueAI struct {
	stats               *configs.GameStats
	dangerEstimator     *estimator.DangerEstimator
	tenpaiProbEstimator *estimator.TenpaiProbEstimator
	noChanges           [4]float64
	logStr              string
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
		estimator.NewDangerEstimator(root),
		estimator.NewTenpaiProbEstimator(stats),
	), nil
}

func NewManueAIWithEstimators(
	stats *configs.GameStats,
	dangerEstimator *estimator.DangerEstimator,
	tenpaiProbEstimator *estimator.TenpaiProbEstimator,
) *ManueAI {
	return &ManueAI{
		stats:               stats,
		dangerEstimator:     dangerEstimator,
		tenpaiProbEstimator: tenpaiProbEstimator,
		noChanges:           [4]float64{},
		logStr:              "",
	}
}

func (a *ManueAI) Initialize() {
	a.logStr = ""
}

func (a *ManueAI) log(str string) {
	fmt.Fprint(os.Stderr, str)
	a.logStr += str
}

func (a *ManueAI) DecideAction(state game.StateAnalyzer, playerID int) (jsontext.Value, error) {
	hc, err := state.HoraCandidate()
	if err != nil {
		return nil, err
	}
	if hc != nil {
		// If it can win, always win
		hora, err := message.NewHora(playerID, hc.Target(), hc.Pai().ToString(), 0, nil, a.logStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create hora message: %w", err)
		}
		res, err := json.Marshal(&hora)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal hora message: %w", err)
		}
		a.logStr = ""
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
	if err != nil {
		return nil, err
	}
	if len(dc) == 0 && len(rdc) == 0 {
		// no action is possible
		return nil, nil
	}

	// my turn

	if state.Players()[playerID].ReachState() == game.Accepted {
		// in reach
		dahai, err := message.NewDahai(playerID, dc[0].ToString(), true, a.logStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create dahai message: %w", err)
		}
		res, err := json.Marshal(&dahai)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal dahai message: %w", err)
		}
		a.logStr = ""
		return res, nil
	}

	ms, err := a.getMetrics(state, playerID, dc, rdc)
	if err != nil {
		return nil, err
	}
	a.printMetrics(ms)
	a.printTenpaiProbs(state, playerID)
	key := a.chooseBestMetric(ms, true)
	fmt.Fprintf(os.Stderr, "decidedKey %s\n", key)
	actionIdx, paiStr, _ := strings.Cut(key, ".")

	reach := actionIdx == "0"
	if reach {
		// reach declaration
		reach, err := message.NewReach(playerID, a.logStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create reach message: %w", err)
		}
		res, err := json.Marshal(&reach)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal reach message: %w", err)
		}
		a.logStr = ""
		return res, nil
	}

	// dahai
	pai, err := game.NewPaiWithName(paiStr)
	if err != nil {
		return nil, err
	}
	isTsumogiri := state.IsTsumoPai(pai)
	dahai, err := message.NewDahai(playerID, paiStr, isTsumogiri, a.logStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create dahai message: %w", err)
	}
	res, err := json.Marshal(&dahai)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dahai message: %w", err)
	}
	a.logStr = ""
	return res, nil
}

func (a *ManueAI) decideFuro(state game.StateAnalyzer, playerID int) (jsontext.Value, error) {
	fc, err := state.FuroCandidates()
	if err != nil {
		return nil, err
	}
	if len(fc) == 0 {
		// no action is possible
		return nil, nil
	}

	// The maximum number of furo candidates is
	// 5 for chi, 1 for pon, and 1 for daiminkan, totaling 7.
	ms := make(metrics, 7)
	player := state.Players()[playerID]

	// Not furo
	noneTehais := player.Tehais()
	noneFuros := player.Furos()
	noneDahai := []game.Pai{*game.Unknown}
	noneMetrics, err := a.getMetricsInternal(state, playerID, noneTehais, noneFuros, noneDahai, false)
	if err != nil {
		return nil, err
	}
	ms["none"] = noneMetrics["none"]

	// Metrics for each furo candidate
	for j, action := range fc {
		tehais := slices.Clone(player.Tehais())
		// remove the consumed tiles from the hand
		for _, pai := range action.Consumed() {
			for i, t := range tehais {
				if t.ID() == pai.ID() {
					tehais = slices.Delete(tehais, i, i+1)
					break
				}
			}
		}
		furos := slices.Clone(player.Furos())
		furos = append(furos, action)
		dahaiCandidates := getUniqueDahais(tehais, func(p game.Pai) bool {
			return isKuikae(action, &p)
		})
		furoMetrics, err := a.getMetricsInternal(state, playerID, tehais, furos, dahaiCandidates, false)
		if err != nil {
			return nil, err
		}
		for k, v := range furoMetrics {
			ms[fmt.Sprintf("%d.%s", j, k)] = v
		}
	}

	a.printMetrics(ms)
	a.printTenpaiProbs(state, playerID)
	key := a.chooseBestMetric(ms, false)
	fmt.Fprintf(os.Stderr, "decidedKey %s\n", key)

	if key == "none" {
		none, err := message.NewSkip(playerID, a.logStr)
		if err != nil {
			return nil, err
		}
		res, err := json.Marshal(&none)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal skip message: %w", err)
		}
		a.logStr = ""
		return res, nil
	}

	actionIdx, _, _ := strings.Cut(key, ".")
	idx, err := strconv.Atoi(actionIdx)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(fc) {
		return nil, fmt.Errorf("invalid furo action index: %s", actionIdx)
	}
	decision := fc[idx]

	target := *decision.Target()
	taken := decision.Taken().ToString()
	switch decision.(type) {
	case *game.Chi:
		var consumed [2]string
		for i, pai := range decision.Consumed() {
			consumed[i] = pai.ToString()
		}
		chiMsg, err := message.NewChi(playerID, target, taken, consumed, a.logStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create chi message: %w", err)
		}
		res, err := json.Marshal(&chiMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal chi message: %w", err)
		}
		a.logStr = ""
		return res, nil
	case *game.Pon:
		var consumed [2]string
		for i, pai := range decision.Consumed() {
			consumed[i] = pai.ToString()
		}
		ponMsg, err := message.NewPon(playerID, target, taken, consumed, a.logStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create pon message: %w", err)
		}
		res, err := json.Marshal(&ponMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal pon message: %w", err)
		}
		a.logStr = ""
		return res, nil
	case *game.Daiminkan:
		var consumed [3]string
		for i, pai := range decision.Consumed() {
			consumed[i] = pai.ToString()
		}
		daiminkanMsg, err := message.NewDaiminkan(playerID, target, taken, consumed, a.logStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create daiminkan message: %w", err)
		}
		res, err := json.Marshal(&daiminkanMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal daiminkan message: %w", err)
		}
		a.logStr = ""
		return res, nil
	default:
		return nil, fmt.Errorf("unknown furo action type: %T", decision)
	}
}

func getUniqueDahais(tehais []game.Pai, del func(game.Pai) bool) []game.Pai {
	unique := game.Pais(slices.Clone(tehais))
	sort.Sort(unique)
	unique = slices.CompactFunc(unique, func(a, b game.Pai) bool {
		return a.ID() == b.ID()
	})
	if del == nil {
		return unique
	}
	unique = slices.DeleteFunc(unique, func(p game.Pai) bool {
		return del(p)
	})
	return unique
}

func isKuikae(furo game.Furo, dahai *game.Pai) bool {
	if dahai.HasSameSymbol(furo.Taken()) {
		return true
	}

	chi, isChi := furo.(*game.Chi)
	if !isChi {
		return false
	}

	pais := chi.Pais()
	number := dahai.Number()
	if number > 3 && number-3 == pais[0].Number() {
		return true
	}
	if number < 7 && number+3 == pais[2].Number() {
		return true
	}
	return false
}
