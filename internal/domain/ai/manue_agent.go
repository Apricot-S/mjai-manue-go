package ai

import (
	"fmt"
	"math/rand/v2"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type ManueAgent struct {
	seed uint64
	rng  *rand.Rand
	deps ManueAgentDeps
}

type ManueAgentDeps struct {
	Stats  ManueStats
	Danger DangerEstimator
}

// ManueStats provides read-only access to immutable statistical data used by
// ManueAgent. Implementations must return stable values for the lifetime of the
// agent after validation.
type ManueStats interface {
	WinScoreStats
	RoundEndStats
	DrawTenpaiStats
	TenpaiEstimatorStats
	RankStats
	DealInStats
}

type WinScoreStats interface {
	NumWins() int
	NumSelfDrawWins() int
	NonDealerWinPointFreqs() map[string]int
	DealerWinPointFreqs() map[string]int
}

type RoundEndStats interface {
	TurnDistribution() []float64
	ExhaustiveDrawRatio() float64
}

type TenpaiEstimatorStats interface {
	YamitenCounts(remainTurns int, numMelds int) (total int, tenpai int, ok bool)
}

type DealInStats interface {
	AvgWinPts() float64
}

type RankStats interface {
	RelativeWinProbs(roundWind wind.Wind, roundNumber int, selfPosition int, otherPosition int) (map[string]float64, bool)
}

type DrawTenpaiStats interface {
	ExhaustiveDrawNotenCount() int
	ExhaustiveDrawTenpaiTurnFreq(turnKey string) (freq int, ok bool)
}

type DangerEstimator interface {
	EstimateDealInProb(state round.StateViewer, self seat.Seat, winner seat.Seat, discard tile.Tile) (float64, error)
}

func NewManueAgent(seed uint64, deps ManueAgentDeps) (*ManueAgent, error) {
	if deps.Stats == nil {
		return nil, fmt.Errorf("cannot create ManueAgent: stats dependency is required")
	}
	if err := validateManueStats(deps.Stats); err != nil {
		return nil, fmt.Errorf("cannot create ManueAgent: %w", err)
	}
	if deps.Danger == nil {
		return nil, fmt.Errorf("cannot create ManueAgent: danger estimator dependency is required")
	}
	agent := &ManueAgent{
		seed: seed,
		deps: deps,
	}
	agent.Reset()
	return agent, nil
}

func (a *ManueAgent) Reset() {
	a.rng = rand.New(rand.NewPCG(a.seed, 0))
}

func (a *ManueAgent) Decide(request Request) (Decision, error) {
	legalActions, err := request.Round.LegalActions(request.Self)
	if err != nil {
		return Decision{}, err
	}
	if len(legalActions) == 0 {
		return Decision{}, fmt.Errorf("cannot decide: no legal actions for player %d", request.Self.Index())
	}

	if win := firstActionOfType[*action.Win](legalActions); win != nil {
		// Always take a winning action when it is legal.
		// The current policy does not allow passing on win opportunities.
		return Decision{Action: win}, nil
	}

	self := request.Round.Player(request.Self)
	if self.CanDiscard() {
		return a.decideSelfTurn(legalActions, request.Round, request.Self, self)
	}
	return a.decideOtherDiscardReaction(legalActions, request.Round, request.Self, self)
}

func (a *ManueAgent) decideSelfTurn(
	legalActions []action.Action,
	state round.StateViewer,
	selfSeat seat.Seat,
	self player.PlayerViewer,
) (Decision, error) {
	if self.RiichiState() == player.RiichiAccepted {
		// After riichi is accepted, always discard the drawn tile.
		// Concealed kan is intentionally ignored even if it is legal.
		if discard := tsumogiriDiscard(legalActions); discard != nil {
			return Decision{Action: discard}, nil
		}
		return Decision{}, fmt.Errorf("cannot decide: no tsumogiri discard after riichi accepted")
	}

	candidates, err := getSelfTurnCandidates(legalActions, self)
	if err != nil {
		return Decision{}, err
	}
	if len(candidates) == 0 {
		return Decision{}, fmt.Errorf("cannot decide self turn: no self-turn candidate")
	}
	if state != nil {
		var err error
		candidates, err = a.evaluateActionCandidates(state, selfSeat, candidates)
		if err != nil {
			return Decision{}, err
		}
	}

	candidate := chooseBestCandidate(candidates, true)
	log := formatCandidateLog(candidates)
	return Decision{
		Action: candidate.action,
		Log:    log,
		Trace:  formatDecisionTrace(log, &candidate),
	}, nil
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

	selfPlayer := state.Player(self)
	context := winEstimateGoalContext{
		melds:          selfPlayer.Melds(),
		roundWind:      state.RoundWind(),
		seatWind:       state.SeatWind(self),
		doraIndicators: state.DoraIndicators(),
		dealer:         state.Dealer() == self,
	}
	goalsByKey, err := scoredWinEstimateGoalsByKey(candidates, context)
	if err != nil {
		return nil, err
	}
	winEstimates, err := winEstimatesFromState(stats, state, self, candidates, goalsByKey, 1000, a.rng)
	if err != nil {
		return nil, err
	}

	exhaustiveDrawProbOnSelfNoWin, err := exhaustiveDrawProbOnSelfNoWin(stats, state.Turn())
	if err != nil {
		return nil, err
	}
	notenTenpaiProb, err := notenExhaustiveDrawTenpaiProb(stats, state.Turn())
	if err != nil {
		return nil, err
	}
	otherWinDists := otherWinScoreDeltaDists(stats, state, self)
	baseTenpaiProbs := currentTenpaiProbs(stats, state, self)

	updated := make([]actionCandidate, len(candidates))
	for i, candidate := range candidates {
		winEstimate, ok := winEstimates[candidate.traceKey]
		if !ok {
			return nil, fmt.Errorf("cannot evaluate candidate %q: missing win estimate", candidate.traceKey)
		}
		selfWinDist := selfWinScoreDeltaDistFromEstimate(self.Index(), state.Dealer().Index(), stats, winEstimate)
		var dealInEstimates []dealInEstimate
		immediateDist := immediateScoreDeltaDist(nil)
		if !candidate.discardTile.IsUnknown() {
			dealInEstimates, err = a.dealInEstimates(state, self, candidate.discardTile)
			if err != nil {
				return nil, fmt.Errorf("cannot evaluate candidate %q deal-in estimates: %w", candidate.traceKey, err)
			}
			immediateDist, err = immediateScoreDeltaDistFromStats(self.Index(), state.Dealer().Index(), dealInEstimates, stats)
			if err != nil {
				return nil, fmt.Errorf("cannot evaluate candidate %q immediate distribution: %w", candidate.traceKey, err)
			}
		}

		tenpaiProbs := baseTenpaiProbs
		tenpaiProbs[self.Index()] = 0
		if candidate.score.shanten <= 0 {
			tenpaiProbs[self.Index()] = 1
		}
		exhaustiveDrawTenpaiProbs := exhaustiveDrawTenpaiProbs(tenpaiProbs, notenTenpaiProb)
		exhaustiveDrawDist := exhaustiveDrawScoreDeltaDistFromTenpaiProbs(exhaustiveDrawTenpaiProbs)
		avgDrawPts := exhaustiveDrawAvgPts(self.Index(), exhaustiveDrawTenpaiProbs)

		score, err := evaluateCandidateFromPreparedDists(
			candidate.score,
			dealInEstimates,
			winEstimate,
			exhaustiveDrawProbOnSelfNoWin,
			avgDrawPts,
			immediateDist,
			selfWinDist,
			exhaustiveDrawDist,
			otherWinDists,
			stats,
			state,
			self,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot evaluate candidate %q: %w", candidate.traceKey, err)
		}
		candidate.score = score
		updated[i] = candidate
	}
	return updated, nil
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

func (a *ManueAgent) decideOtherDiscardReaction(
	legalActions []action.Action,
	state round.StateViewer,
	selfSeat seat.Seat,
	self player.PlayerViewer,
) (Decision, error) {
	if state == nil {
		return Decision{}, fmt.Errorf("cannot decide other discard reaction: state is required")
	}
	if self == nil {
		return Decision{}, fmt.Errorf("cannot decide other discard reaction: self player is required")
	}

	candidates, err := getOtherDiscardReactionCandidates(legalActions, self)
	if err != nil {
		return Decision{}, err
	}
	if len(candidates) == 0 {
		return Decision{}, fmt.Errorf("cannot decide other discard reaction: no reaction candidate")
	}

	candidates, err = a.evaluateActionCandidates(state, selfSeat, candidates)
	if err != nil {
		return Decision{}, err
	}

	candidate := chooseBestCandidate(candidates, false)
	log := formatCandidateLog(candidates)
	return Decision{
		Action: candidate.action,
		Log:    log,
		Trace:  formatDecisionTrace(log, &candidate),
	}, nil
}

func firstActionOfType[T action.Action](actions []action.Action) T {
	for _, a := range actions {
		if typed, ok := a.(T); ok {
			return typed
		}
	}
	var zero T
	return zero
}

func tsumogiriDiscard(actions []action.Action) *action.Discard {
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if ok && discard.Tsumogiri() {
			return discard
		}
	}
	return nil
}
