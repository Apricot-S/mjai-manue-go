package ai

import (
	"fmt"
	"math/rand/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type actionEvaluationContext struct {
	stats                         ManueStats
	state                         round.StateViewer
	self                          seat.Seat
	winEstimates                  map[string]winEstimate
	exhaustiveDrawProbOnSelfNoWin float64
	notenTenpaiProb               float64
	otherWinDists                 []scoreDeltaProbDist
	baseTenpaiProbs               [common.NumPlayers]float64
}

func (a *ManueAgent) evaluateActionCandidates(
	state round.StateViewer,
	self seat.Seat,
	candidates []actionCandidate,
) ([]actionCandidate, error) {
	stats := a.deps.Stats
	if stats == nil {
		return candidates, nil
	}

	context, err := newActionEvaluationContext(stats, state, self, candidates, a.rng)
	if err != nil {
		return nil, err
	}

	updated := make([]actionCandidate, len(candidates))
	for i, candidate := range candidates {
		updatedCandidate, err := a.evaluateActionCandidate(context, candidate)
		if err != nil {
			return nil, fmt.Errorf("cannot evaluate candidate %q: %w", candidate.traceKey, err)
		}
		updated[i] = updatedCandidate
	}
	return updated, nil
}

func newActionEvaluationContext(
	stats ManueStats,
	state round.StateViewer,
	self seat.Seat,
	candidates []actionCandidate,
	rng *rand.Rand,
) (actionEvaluationContext, error) {
	selfPlayer := state.Player(self)
	goalContext := winEstimateGoalContext{
		melds:          selfPlayer.Melds(),
		roundWind:      state.RoundWind(),
		seatWind:       state.SeatWind(self),
		doraIndicators: state.DoraIndicators(),
		dealer:         state.Dealer() == self,
	}
	goalsByKey, err := scoredWinEstimateGoalsByKey(candidates, goalContext)
	if err != nil {
		return actionEvaluationContext{}, err
	}
	winEstimates, err := winEstimatesFromState(stats, state, self, candidates, goalsByKey, 1000, rng)
	if err != nil {
		return actionEvaluationContext{}, err
	}

	exhaustiveDrawProbOnSelfNoWin, err := exhaustiveDrawProbOnSelfNoWin(stats, state.Turn())
	if err != nil {
		return actionEvaluationContext{}, err
	}
	notenTenpaiProb, err := notenExhaustiveDrawTenpaiProb(stats, state.Turn())
	if err != nil {
		return actionEvaluationContext{}, err
	}

	return actionEvaluationContext{
		stats:                         stats,
		state:                         state,
		self:                          self,
		winEstimates:                  winEstimates,
		exhaustiveDrawProbOnSelfNoWin: exhaustiveDrawProbOnSelfNoWin,
		notenTenpaiProb:               notenTenpaiProb,
		otherWinDists:                 otherWinScoreDeltaDists(stats, state, self),
		baseTenpaiProbs:               currentTenpaiProbs(stats, state, self),
	}, nil
}

func (a *ManueAgent) evaluateActionCandidate(
	context actionEvaluationContext,
	candidate actionCandidate,
) (actionCandidate, error) {
	winEstimate, ok := context.winEstimates[candidate.traceKey]
	if !ok {
		return actionCandidate{}, fmt.Errorf("missing win estimate")
	}
	selfWinDist := selfWinScoreDeltaDistFromEstimate(
		context.self.Index(),
		context.state.Dealer().Index(),
		context.stats,
		winEstimate,
	)

	dealInEstimates, immediateDist, err := a.immediateDealInEvaluation(context, candidate)
	if err != nil {
		return actionCandidate{}, err
	}

	exhaustiveDrawDist, exhaustiveDrawAveragePoints := candidateExhaustiveDrawEvaluation(context, candidate)
	score, err := evaluateCandidateFromComponents(
		candidate.score,
		dealInEstimates,
		winEstimate,
		context.exhaustiveDrawProbOnSelfNoWin,
		exhaustiveDrawAveragePoints,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		context.otherWinDists,
		context.stats,
		context.state,
		context.self,
	)
	if err != nil {
		return actionCandidate{}, err
	}
	candidate.score = score
	return candidate, nil
}

func (a *ManueAgent) immediateDealInEvaluation(
	context actionEvaluationContext,
	candidate actionCandidate,
) ([]dealInEstimate, scoreDeltaProbDist, error) {
	if candidate.discardTile.IsUnknown() {
		return nil, immediateScoreDeltaDist(nil), nil
	}

	dealInEstimates, err := a.dealInEstimates(context.state, context.self, candidate.discardTile)
	if err != nil {
		return nil, scoreDeltaProbDist{}, fmt.Errorf("deal-in estimates: %w", err)
	}
	immediateDist, err := immediateScoreDeltaDistFromStats(
		context.self.Index(),
		context.state.Dealer().Index(),
		dealInEstimates,
		context.stats,
	)
	if err != nil {
		return nil, scoreDeltaProbDist{}, fmt.Errorf("immediate distribution: %w", err)
	}
	return dealInEstimates, immediateDist, nil
}

func candidateExhaustiveDrawEvaluation(
	context actionEvaluationContext,
	candidate actionCandidate,
) (scoreDeltaProbDist, float64) {
	tenpaiProbs := context.baseTenpaiProbs
	tenpaiProbs[context.self.Index()] = 0
	if candidate.score.shanten <= 0 {
		tenpaiProbs[context.self.Index()] = 1
	}
	exhaustiveDrawTenpaiProbs := exhaustiveDrawTenpaiProbs(tenpaiProbs, context.notenTenpaiProb)
	return exhaustiveDrawScoreDeltaDistFromTenpaiProbs(exhaustiveDrawTenpaiProbs),
		exhaustiveDrawAvgPts(context.self.Index(), exhaustiveDrawTenpaiProbs)
}

func (a *ManueAgent) dealInEstimates(
	state round.StateViewer,
	self seat.Seat,
	discard tile.Tile,
) ([]dealInEstimate, error) {
	if a.deps.Danger == nil {
		return nil, nil
	}
	remainTurns := stateNumRemainTurns(state)
	estimates := make([]dealInEstimate, 0, common.NumPlayers-1)
	for i := range common.NumPlayers {
		winner := seat.MustSeat(i)
		if winner == self {
			continue
		}
		winnerPlayer := state.Player(winner)
		tenpai := tenpaiProb(a.deps.Stats, winnerPlayer.RiichiState() != player.NotRiichi, remainTurns, len(winnerPlayer.Melds()))
		rawProb, err := a.deps.Danger.EstimateDealInProb(state, self, winner, discard)
		if err != nil {
			return nil, err
		}
		estimates = append(estimates, dealInEstimate{
			winnerID: winner.Index(),
			prob:     tenpai * rawProb,
		})
	}
	return estimates, nil
}

func currentTenpaiProbs(stats TenpaiEstimatorStats, state round.StateViewer, self seat.Seat) [4]float64 {
	remainTurns := stateNumRemainTurns(state)
	var probs [4]float64
	for i := range common.NumPlayers {
		playerSeat := seat.MustSeat(i)
		if playerSeat == self {
			continue
		}
		p := state.Player(playerSeat)
		probs[i] = tenpaiProb(stats, p.RiichiState() != player.NotRiichi, remainTurns, len(p.Melds()))
	}
	return probs
}

func otherWinScoreDeltaDists(stats WinScoreStats, state round.StateViewer, self seat.Seat) []scoreDeltaProbDist {
	dists := make([]scoreDeltaProbDist, 0, common.NumPlayers-1)
	for i := range common.NumPlayers {
		actor := seat.MustSeat(i)
		if actor == self {
			continue
		}
		dists = append(dists, randomWinScoreDeltaDistFromStats(actor.Index(), state.Dealer().Index(), stats))
	}
	return dists
}

func stateNumRemainTurns(state interface{ NumLeftTiles() int }) int {
	return state.NumLeftTiles() / common.NumPlayers
}
