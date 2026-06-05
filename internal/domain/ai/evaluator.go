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

const defaultWinEstimateTrials = 1000

type candidateEvaluationContext struct {
	stats                         ManueStats
	state                         round.StateViewer
	self                          seat.Seat
	winEstimates                  map[string]winEstimate
	winEstimateGoalCounts         []int
	exhaustiveDrawProbOnSelfNoWin float64
	notenTenpaiProb               float64
	otherWinDists                 []scoreDeltaProbDist
	baseTenpaiProbs               [common.NumPlayers]float64
}

type candidateEvaluationSummary struct {
	winEstimateGoalCounts []int
}

type candidateEvaluator struct {
	stats  ManueStats
	danger DangerEstimator
	rng    *rand.Rand
	trials int
}

func newCandidateEvaluator(stats ManueStats, danger DangerEstimator, rng *rand.Rand) candidateEvaluator {
	return candidateEvaluator{
		stats:  stats,
		danger: danger,
		rng:    rng,
		trials: defaultWinEstimateTrials,
	}
}

func (e candidateEvaluator) evaluateCandidates(
	state round.StateViewer,
	self seat.Seat,
	candidates []actionCandidate,
) ([]actionCandidate, candidateEvaluationSummary, error) {
	context, err := e.newEvaluationContext(state, self, candidates)
	if err != nil {
		return nil, candidateEvaluationSummary{}, err
	}

	updated := make([]actionCandidate, len(candidates))
	for i, candidate := range candidates {
		updatedCandidate, err := e.evaluateCandidate(context, candidate)
		if err != nil {
			return nil, candidateEvaluationSummary{}, fmt.Errorf("cannot evaluate candidate %q: %w", candidate.traceKey, err)
		}
		updated[i] = updatedCandidate
	}
	return updated, candidateEvaluationSummary{
		winEstimateGoalCounts: context.winEstimateGoalCounts,
	}, nil
}

func (e candidateEvaluator) newEvaluationContext(
	state round.StateViewer,
	self seat.Seat,
	candidates []actionCandidate,
) (candidateEvaluationContext, error) {
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
		return candidateEvaluationContext{}, err
	}
	winEstimateGoalCounts := countWinEstimateGoalsByGroup(candidates, goalsByKey)
	winEstimates, err := winEstimatesFromState(e.stats, state, self, candidates, goalsByKey, e.trials, e.rng)
	if err != nil {
		return candidateEvaluationContext{}, err
	}

	exhaustiveDrawProbOnSelfNoWin, err := exhaustiveDrawProbOnSelfNoWin(e.stats, state.Turn())
	if err != nil {
		return candidateEvaluationContext{}, err
	}
	notenTenpaiProb, err := notenExhaustiveDrawTenpaiProb(e.stats, state.Turn())
	if err != nil {
		return candidateEvaluationContext{}, err
	}

	return candidateEvaluationContext{
		stats:                         e.stats,
		state:                         state,
		self:                          self,
		winEstimates:                  winEstimates,
		winEstimateGoalCounts:         winEstimateGoalCounts,
		exhaustiveDrawProbOnSelfNoWin: exhaustiveDrawProbOnSelfNoWin,
		notenTenpaiProb:               notenTenpaiProb,
		otherWinDists:                 otherWinScoreDeltaDists(e.stats, state, self),
		baseTenpaiProbs:               currentTenpaiProbs(e.stats, state, self),
	}, nil
}

func (e candidateEvaluator) evaluateCandidate(
	context candidateEvaluationContext,
	candidate actionCandidate,
) (actionCandidate, error) {
	winEstimate, ok := context.winEstimates[candidate.traceKey]
	if !ok {
		return actionCandidate{}, fmt.Errorf("missing win estimate")
	}
	selfWinDist := winScoreDeltaDistFromPointsDist(
		context.self.Index(),
		context.state.Dealer().Index(),
		context.stats,
		winEstimate.pointsDist,
	)

	dealInEstimates, immediateDist, err := e.immediateDealInEvaluation(context, candidate)
	if err != nil {
		return actionCandidate{}, err
	}

	exhaustiveDrawDist, exhaustiveDrawAveragePoints := candidateExhaustiveDrawEvaluation(context, candidate)
	score, err := evaluateCandidateFromComponents(
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

func (e candidateEvaluator) immediateDealInEvaluation(
	context candidateEvaluationContext,
	candidate actionCandidate,
) ([]dealInEstimate, scoreDeltaProbDist, error) {
	if candidate.discardTile.IsUnknown() {
		return nil, immediateScoreDeltaDist(nil), nil
	}

	dealInEstimates, err := e.dealInEstimates(context.state, context.self, candidate.discardTile)
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
	context candidateEvaluationContext,
	candidate actionCandidate,
) (scoreDeltaProbDist, float64) {
	tenpaiProbs := context.baseTenpaiProbs
	tenpaiProbs[context.self.Index()] = 0
	if candidate.shanten <= 0 {
		tenpaiProbs[context.self.Index()] = 1
	}
	exhaustiveDrawTenpaiProbs := exhaustiveDrawTenpaiProbs(tenpaiProbs, context.notenTenpaiProb)
	dist := exhaustiveDrawScoreDeltaDist(exhaustiveDrawTenpaiProbs)
	return dist, dist.expected()[context.self.Index()]
}

func (e candidateEvaluator) dealInEstimates(
	state round.StateViewer,
	self seat.Seat,
	discard tile.Tile,
) ([]dealInEstimate, error) {
	if e.danger == nil {
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
		tenpai := tenpaiProb(e.stats, winnerPlayer.RiichiState() != player.NotRiichi, remainTurns, len(winnerPlayer.Melds()))
		rawProb, err := e.danger.EstimateDealInProb(state, self, winner, discard)
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

func currentTenpaiProbs(stats TenpaiEstimatorStats, state round.StateViewer, self seat.Seat) [common.NumPlayers]float64 {
	remainTurns := stateNumRemainTurns(state)
	var probs [common.NumPlayers]float64
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
