package ai

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type candidateScore struct {
	// avgRank is the average rank.
	avgRank float64
	// expPts is the expected points.
	expPts float64
	// dealInProb is the deal-in probability.
	dealInProb float64
	// winProb is the win probability.
	winProb float64
	// drawProb is the draw probability.
	drawProb float64
	// othersWinProb is the other players' win probability.
	othersWinProb float64
	// avgWinPts is the average win points.
	avgWinPts float64
	// avgDrawPts is the average draw points.
	avgDrawPts float64
	// shanten is the shanten number.
	shanten int
	// red indicates whether the candidate discards a red tile.
	red bool
}

type actionCandidate struct {
	traceKey         string
	action           action.Action
	riichi           bool
	discardTile      tile.Tile
	turnHand         *hand.VisibleHand
	afterDiscardHand *hand.VisibleHand
	shantenGoals     []service.Goal
	score            candidateScore
}

type winEstimateGoal struct {
	service.Goal
	points float64
}

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

type candidateScoreDeltaDists struct {
	immediateDist      scoreDeltaProbDist
	selfWinDist        scoreDeltaProbDist
	exhaustiveDrawDist scoreDeltaProbDist
	otherWinDists      []scoreDeltaProbDist
}

func chooseBestCandidate(candidates []actionCandidate, preferBlack bool) actionCandidate {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if compareCandidate(candidate, best, preferBlack) < 0 {
			best = candidate
		}
	}
	return best
}

func compareCandidate(lhs, rhs actionCandidate, preferBlack bool) int {
	if r := compareCandidateScore(&lhs.score, &rhs.score, preferBlack); r != 0 {
		return r
	}
	return compareCandidateFallback(lhs, rhs)
}

// compareCandidateFallback preserves Phase 1 behavior until score calculation is migrated.
func compareCandidateFallback(lhs, rhs actionCandidate) int {
	if lhs.riichi && !rhs.riichi {
		return -1
	}
	if !lhs.riichi && rhs.riichi {
		return 1
	}
	return 0
}

func compareCandidateScore(lhs, rhs *candidateScore, preferBlack bool) int {
	if lhs.avgRank < rhs.avgRank {
		return -1
	}
	if lhs.avgRank > rhs.avgRank {
		return 1
	}
	if lhs.expPts > rhs.expPts {
		return -1
	}
	if lhs.expPts < rhs.expPts {
		return 1
	}
	if preferBlack {
		if !lhs.red && rhs.red {
			return -1
		}
		if lhs.red && !rhs.red {
			return 1
		}
	}
	return 0
}

// evaluateCandidateScore fills the fields derived from the final score-change
// distribution while preserving candidate-local estimates such as probabilities
// and shanten.
func evaluateCandidateScore(
	score candidateScore,
	scoreChanges scoreDeltaProbDist,
	selfID int,
	selfScore float64,
	selfPosition int,
	opponents []rankOpponent,
) candidateScore {
	score.expPts = expectedPts(selfID, scoreChanges)
	score.avgRank = averageRank(scoreChanges, selfID, selfScore, selfPosition, opponents)
	return score
}

func evaluateCandidateScoreFromState(
	score candidateScore,
	scoreChanges scoreDeltaProbDist,
	stats RankStats,
	state rankStateViewer,
	self seat.Seat,
) candidateScore {
	scores := state.Scores()
	startingDealer := state.StartingDealer()
	return evaluateCandidateScore(
		score,
		scoreChanges,
		self.Index(),
		float64(scores[self.Index()]),
		self.DistanceFrom(startingDealer),
		buildRankOpponents(stats, state, self),
	)
}

func applyDealInEstimatesToCandidateScore(score candidateScore, estimates []dealInEstimate) (candidateScore, float64, error) {
	safeProb, err := safeProb(estimates)
	if err != nil {
		return candidateScore{}, 0, err
	}
	score.dealInProb = 1.0 - safeProb
	return score, safeProb, nil
}

func applyWinEstimateToCandidateScore(score candidateScore, estimate winEstimate) (candidateScore, error) {
	if estimate.prob < 0.0 || estimate.prob > 1.0 {
		return candidateScore{}, fmt.Errorf("cannot apply win estimate: win probability must be between 0 and 1")
	}
	score.winProb = estimate.prob
	score.avgWinPts = estimate.avgPts
	return score, nil
}

func applyWinEstimatesToCandidates(
	candidates []actionCandidate,
	estimates map[string]winEstimate,
) ([]actionCandidate, error) {
	updated := make([]actionCandidate, len(candidates))
	for i, candidate := range candidates {
		estimate, ok := estimates[candidate.traceKey]
		if !ok {
			return nil, fmt.Errorf("cannot apply win estimates: missing estimate for %q", candidate.traceKey)
		}
		score, err := applyWinEstimateToCandidateScore(candidate.score, estimate)
		if err != nil {
			return nil, fmt.Errorf("cannot apply win estimate for %q: %w", candidate.traceKey, err)
		}
		candidate.score = score
		updated[i] = candidate
	}
	return updated, nil
}

func applyRoundEndProbsToCandidateScore(score candidateScore, exhaustiveDrawProbOnSelfNoWin float64) (candidateScore, error) {
	drawProb, othersWinProb, err := remainingRoundEndProbs(score.winProb, exhaustiveDrawProbOnSelfNoWin)
	if err != nil {
		return candidateScore{}, err
	}
	score.drawProb = drawProb
	score.othersWinProb = othersWinProb
	return score, nil
}

func applyAveragePtsToCandidateScore(score candidateScore, avgWinPts float64, avgDrawPts float64) candidateScore {
	score.avgWinPts = avgWinPts
	score.avgDrawPts = avgDrawPts
	return score
}

func candidateTotalScoreDeltaDist(
	score candidateScore,
	immediateDist scoreDeltaProbDist,
	selfWinDist scoreDeltaProbDist,
	exhaustiveDrawDist scoreDeltaProbDist,
	otherWinDists []scoreDeltaProbDist,
) scoreDeltaProbDist {
	futureDist := futureScoreDeltaDist(
		selfWinDist,
		score.winProb,
		exhaustiveDrawDist,
		score.drawProb,
		otherWinDists,
		score.othersWinProb,
	)
	return totalScoreDeltaDist(immediateDist, futureDist)
}

func evaluateCandidateFromPreparedDists(
	score candidateScore,
	dealInEstimates []dealInEstimate,
	winEstimate winEstimate,
	exhaustiveDrawProbOnSelfNoWin float64,
	avgDrawPts float64,
	dists candidateScoreDeltaDists,
	rankStats RankStats,
	state rankStateViewer,
	self seat.Seat,
) (candidateScore, error) {
	score, _, err := applyDealInEstimatesToCandidateScore(score, dealInEstimates)
	if err != nil {
		return candidateScore{}, err
	}
	score, err = applyWinEstimateToCandidateScore(score, winEstimate)
	if err != nil {
		return candidateScore{}, err
	}
	score, err = applyRoundEndProbsToCandidateScore(score, exhaustiveDrawProbOnSelfNoWin)
	if err != nil {
		return candidateScore{}, err
	}
	score = applyAveragePtsToCandidateScore(score, score.avgWinPts, avgDrawPts)
	scoreChanges := candidateTotalScoreDeltaDist(
		score,
		dists.immediateDist,
		dists.selfWinDist,
		dists.exhaustiveDrawDist,
		dists.otherWinDists,
	)
	return evaluateCandidateScoreFromState(score, scoreChanges, rankStats, state, self), nil
}

func candidateTraceKeys(candidates []actionCandidate) ([]string, error) {
	keys := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if candidate.traceKey == "" {
			return nil, fmt.Errorf("cannot build candidate keys: trace key must not be empty")
		}
		if _, ok := seen[candidate.traceKey]; ok {
			return nil, fmt.Errorf("cannot build candidate keys: duplicate trace key %q", candidate.traceKey)
		}
		seen[candidate.traceKey] = struct{}{}
		keys = append(keys, candidate.traceKey)
	}
	return keys, nil
}

func winEstimatesForCandidates(candidates []actionCandidate, trials []map[string]float64) (map[string]winEstimate, error) {
	keys, err := candidateTraceKeys(candidates)
	if err != nil {
		return nil, err
	}
	return winEstimatesFromTrials(keys, trials)
}

func applyWinEstimateTrialsToCandidates(
	candidates []actionCandidate,
	trials []map[string]float64,
) ([]actionCandidate, error) {
	estimates, err := winEstimatesForCandidates(candidates, trials)
	if err != nil {
		return nil, err
	}
	return applyWinEstimatesToCandidates(candidates, estimates)
}

func filteredWinEstimateGoals(candidate actionCandidate) []service.Goal {
	goals := make([]service.Goal, 0, len(candidate.shantenGoals))
	for _, goal := range candidate.shantenGoals {
		if candidate.riichi && goal.Shanten > 0 {
			continue
		}
		if candidate.score.shanten > 3 && goal.Shanten > candidate.score.shanten {
			continue
		}
		discardID := candidate.discardTile.RemoveRed().ID()
		if goal.ThrowableVector[discardID] <= 0 {
			continue
		}
		goals = append(goals, goal)
	}
	return goals
}

func scoredWinEstimateGoals(candidate actionCandidate, context winEstimateGoalContext) ([]winEstimateGoal, error) {
	if candidate.turnHand == nil {
		return nil, fmt.Errorf("cannot score win estimate goals: turn hand must not be nil")
	}

	goals := filteredWinEstimateGoals(candidate)
	scoredGoals := make([]winEstimateGoal, 0, len(goals))
	for _, goal := range goals {
		fu, han, _ := service.CalculateFuHan(
			candidate.turnHand,
			goal.Blocks,
			context.melds,
			context.roundWind,
			context.seatWind,
			context.doraIndicators,
			candidate.riichi,
		)
		points := service.RonPoints(fu, han, context.dealer)
		if points <= 0 {
			continue
		}
		scoredGoals = append(scoredGoals, winEstimateGoal{
			Goal:   goal,
			points: float64(points),
		})
	}
	return scoredGoals, nil
}

func scoredWinEstimateGoalsByKey(
	candidates []actionCandidate,
	context winEstimateGoalContext,
) (map[string][]winEstimateGoal, error) {
	if _, err := candidateTraceKeys(candidates); err != nil {
		return nil, err
	}

	goalsByKey := make(map[string][]winEstimateGoal, len(candidates))
	for _, candidate := range candidates {
		goals, err := scoredWinEstimateGoals(candidate, context)
		if err != nil {
			return nil, fmt.Errorf("cannot score win estimate goals for %q: %w", candidate.traceKey, err)
		}
		goalsByKey[candidate.traceKey] = goals
	}
	return goalsByKey, nil
}

func trialTileCounts(tiles []tile.Tile) hand.TileCounts34 {
	var counts hand.TileCounts34
	for _, t := range tiles {
		counts[t.RemoveRed().ID()]++
	}
	return counts
}

func wallTilesFromCounts(counts hand.TileCounts34) ([]tile.Tile, error) {
	for id, count := range counts {
		if count < 0 {
			return nil, fmt.Errorf("cannot build wall tiles: tile %s count must be non-negative", tile.MustTileFromID(id))
		}
	}
	return (&counts).ToTiles(), nil
}

func unseenWallFromVisibleTiles(visibleTiles []tile.Tile) ([]tile.Tile, error) {
	var counts hand.TileCounts34
	for id := range counts {
		counts[id] = 4
	}
	for _, visible := range visibleTiles {
		if visible.IsUnknown() {
			return nil, fmt.Errorf("cannot build unseen wall: visible tile must not be unknown")
		}
		id := visible.RemoveRed().ID()
		counts[id]--
		if counts[id] < 0 {
			return nil, fmt.Errorf("cannot build unseen wall: tile %s is visible more than 4 times", tile.MustTileFromID(id))
		}
	}
	return wallTilesFromCounts(counts)
}

func trialTilesFromWall(wall []tile.Tile, numDraws int) ([]tile.Tile, error) {
	if numDraws < 0 {
		return nil, fmt.Errorf("cannot build trial tiles: numDraws must be non-negative")
	}
	if numDraws > len(wall) {
		return nil, fmt.Errorf("cannot build trial tiles: numDraws %d exceeds wall length %d", numDraws, len(wall))
	}
	return slices.Clone(wall[:numDraws]), nil
}

func canAchieveGoalWithTrialTiles(goal service.Goal, trialTiles hand.TileCounts34) bool {
	for id, required := range goal.RequiredVector {
		if required > trialTiles[id] {
			return false
		}
	}
	return true
}

func trialWinPts(goals []winEstimateGoal, trialTiles hand.TileCounts34) (float64, bool, error) {
	best := 0.0
	for _, goal := range goals {
		if goal.points <= 0 {
			return 0, false, fmt.Errorf("cannot calculate trial win points: goal points must be positive")
		}
		if !canAchieveGoalWithTrialTiles(goal.Goal, trialTiles) {
			continue
		}
		best = max(best, goal.points)
	}
	return best, best > 0, nil
}

func candidateTrialWinPts(
	candidates []actionCandidate,
	goalsByKey map[string][]winEstimateGoal,
	trialTiles hand.TileCounts34,
) (map[string]float64, error) {
	if _, err := candidateTraceKeys(candidates); err != nil {
		return nil, err
	}

	winPtsByKey := make(map[string]float64)
	for _, candidate := range candidates {
		goals, ok := goalsByKey[candidate.traceKey]
		if !ok {
			return nil, fmt.Errorf("cannot calculate candidate trial win points: missing goals for %q", candidate.traceKey)
		}
		points, won, err := trialWinPts(goals, trialTiles)
		if err != nil {
			return nil, fmt.Errorf("cannot calculate candidate trial win points for %q: %w", candidate.traceKey, err)
		}
		if won {
			winPtsByKey[candidate.traceKey] = points
		}
	}
	return winPtsByKey, nil
}

func winEstimatesFromTrialTiles(
	candidates []actionCandidate,
	goalsByKey map[string][]winEstimateGoal,
	trials [][]tile.Tile,
) (map[string]winEstimate, error) {
	trialResults := make([]map[string]float64, 0, len(trials))
	for i, trial := range trials {
		points, err := candidateTrialWinPts(candidates, goalsByKey, trialTileCounts(trial))
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimates from trial %d: %w", i, err)
		}
		trialResults = append(trialResults, points)
	}
	return winEstimatesForCandidates(candidates, trialResults)
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
		shuffled := slices.Clone(wall)
		rng.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		trialTiles, err := trialTilesFromWall(shuffled, numDraws)
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimates from shuffled wall trial %d: %w", i, err)
		}
		points, err := candidateTrialWinPts(candidates, goalsByKey, trialTileCounts(trialTiles))
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

func getSelfTurnCandidates(actions []action.Action, self player.PlayerViewer) ([]actionCandidate, error) {
	h, err := selfTurnHand(self)
	if err != nil {
		return nil, fmt.Errorf("cannot build self-turn candidates: %w", err)
	}
	_, turnGoals := service.AnalyzeShanten(h, service.AllowedExtraTiles(1))

	// Self-turn candidates currently cover discard and riichi+discard only.
	// Concealed kan, promoted kan, and kyushukyuhai are intentionally not
	// selected.
	riichi := firstActionOfType[*action.Riichi](actions)
	var candidates []actionCandidate
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if !ok {
			continue
		}

		afterDiscard, err := h.Discard(discard.Tile())
		if err != nil {
			return nil, fmt.Errorf("cannot build self-turn candidate for %s: %w", discard.Tile(), err)
		}
		shanten, _ := service.AnalyzeShanten(afterDiscard, service.AllowedExtraTiles(1))
		if riichi != nil && shanten <= 0 {
			// Match Manue's riichi candidate filtering: only regular-hand shanten
			// from AnalyzeShanten is considered here, so Seven Pairs and
			// Thirteen Orphans do not create riichi candidates.
			candidates = append(candidates, buildSelfTurnCandidate(riichi, discard.Tile(), h, afterDiscard, shanten, turnGoals, true))
		}
		candidates = append(candidates, buildSelfTurnCandidate(discard, discard.Tile(), h, afterDiscard, shanten, turnGoals, false))
	}
	return candidates, nil
}

func selfTurnHand(self player.PlayerViewer) (*hand.VisibleHand, error) {
	h, ok := self.Hand()
	if !ok {
		return nil, fmt.Errorf("self hand is not visible")
	}
	drawnTile := self.DrawnTile()
	if drawnTile == nil {
		return h, nil
	}
	withDrawnTile, err := h.Draw(*drawnTile)
	if err != nil {
		return nil, fmt.Errorf("cannot add drawn tile %s to self hand: %w", *drawnTile, err)
	}
	return withDrawnTile, nil
}

func buildSelfTurnCandidate(
	immediateAction action.Action,
	discardTile tile.Tile,
	turnHand *hand.VisibleHand,
	afterDiscardHand *hand.VisibleHand,
	shanten int,
	goals []service.Goal,
	riichi bool,
) actionCandidate {
	return actionCandidate{
		traceKey:         formatDiscardTraceKey(riichi, discardTile),
		action:           immediateAction,
		riichi:           riichi,
		discardTile:      discardTile,
		turnHand:         turnHand,
		afterDiscardHand: afterDiscardHand,
		shantenGoals:     goals,
		score:            scoreDiscardCandidate(discardTile, shanten),
	}
}

func scoreDiscardCandidate(discardTile tile.Tile, shanten int) candidateScore {
	return candidateScore{
		// Phase 2 scaffold: the full expected-value fields are filled by later
		// configs/estimator migration. Keep all choices tied except red fallback.
		avgRank: 0,
		expPts:  0,
		shanten: shanten,
		red:     discardTile.IsRed(),
	}
}
