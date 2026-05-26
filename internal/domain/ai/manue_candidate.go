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
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service/block"
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
	scoreAsRiichi    bool
	discardTile      tile.Tile
	melds            []meld.Meld
	turnHand         *hand.VisibleHand
	afterDiscardHand *hand.VisibleHand
	baseShanten      int
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

func chooseBestCandidate(candidates []actionCandidate, preferBlack bool) actionCandidate {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if compareCandidateScore(&candidate.score, &best.score, preferBlack) < 0 {
			best = candidate
		}
	}
	return best
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
	immediateDist scoreDeltaProbDist,
	selfWinDist scoreDeltaProbDist,
	exhaustiveDrawDist scoreDeltaProbDist,
	otherWinDists []scoreDeltaProbDist,
	rankStats RankStats,
	state rankStateViewer,
	self seat.Seat,
) (candidateScore, error) {
	safeProb, err := safeProb(dealInEstimates)
	if err != nil {
		return candidateScore{}, err
	}
	score.dealInProb = 1.0 - safeProb
	if winEstimate.prob < 0.0 || winEstimate.prob > 1.0 {
		return candidateScore{}, fmt.Errorf("cannot evaluate candidate: win probability must be between 0 and 1")
	}
	score.winProb = winEstimate.prob
	score.avgWinPts = winEstimate.avgPts
	drawProb, othersWinProb, err := remainingRoundEndProbs(score.winProb, exhaustiveDrawProbOnSelfNoWin)
	if err != nil {
		return candidateScore{}, err
	}
	score.drawProb = drawProb
	score.othersWinProb = othersWinProb
	score.avgDrawPts = avgDrawPts
	scoreChanges := candidateTotalScoreDeltaDist(
		score,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		otherWinDists,
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

func filteredWinEstimateGoals(candidate actionCandidate) []service.Goal {
	goals := make([]service.Goal, 0, len(candidate.shantenGoals))
	for _, goal := range candidate.shantenGoals {
		if candidate.scoreAsRiichi && goal.Shanten > 0 {
			continue
		}
		if candidate.baseShanten > 3 && goal.Shanten > candidate.baseShanten {
			continue
		}
		if !candidate.discardTile.IsUnknown() {
			discardID := candidate.discardTile.RemoveRed().ID()
			if goal.ThrowableVector[discardID] <= 0 {
				continue
			}
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
		scoringHand, err := scoringHandForGoal(candidate.turnHand, goal.Blocks)
		if err != nil {
			return nil, err
		}
		fu, han, _ := service.CalculateFuHan(
			scoringHand,
			goal.Blocks,
			context.melds,
			context.roundWind,
			context.seatWind,
			context.doraIndicators,
			candidate.scoreAsRiichi,
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

func scoringHandForGoal(sourceHand *hand.VisibleHand, blocks []block.Block) (*hand.VisibleHand, error) {
	redCounts := make(map[int]int, 3)
	for _, t := range sourceHand.ToTiles() {
		if t.IsRed() {
			redCounts[t.RemoveRed().ID()]++
		}
	}

	tiles := make([]tile.Tile, 0, 14)
	for _, b := range blocks {
		for _, t := range b.ToTiles() {
			normal := t.RemoveRed()
			if redCounts[normal.ID()] > 0 {
				tiles = append(tiles, normal.AddRed())
				redCounts[normal.ID()]--
				continue
			}
			tiles = append(tiles, normal)
		}
	}

	scoringHand, err := hand.NewVisibleHand(tiles)
	if err != nil {
		return nil, fmt.Errorf("cannot build scoring hand for win estimate goal: %w", err)
	}
	return scoringHand, nil
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
		candidateContext := context
		if candidate.melds != nil {
			candidateContext.melds = candidate.melds
		}
		goals, err := scoredWinEstimateGoals(candidate, candidateContext)
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

func candidateShanten(discardTile tile.Tile, baseShanten int, goals []service.Goal) int {
	if discardTile.IsUnknown() {
		return baseShanten
	}
	discardID := discardTile.RemoveRed().ID()
	shanten := service.InfinityShanten
	for _, goal := range goals {
		if goal.ThrowableVector[discardID] > 0 && goal.Shanten < shanten {
			shanten = goal.Shanten
		}
	}
	return shanten
}

func getSelfTurnCandidates(actions []action.Action, self player.PlayerViewer) ([]actionCandidate, error) {
	h, err := selfTurnHand(self)
	if err != nil {
		return nil, fmt.Errorf("cannot build self-turn candidates: %w", err)
	}
	turnShanten, turnGoals := service.AnalyzeShanten(h, service.AllowedExtraTiles(1))

	// Self-turn candidates currently cover discard and riichi+discard only.
	// Concealed kan, promoted kan, and kyushukyuhai are intentionally not
	// selected.
	riichi := firstActionOfType[*action.Riichi](actions)
	riichiDeclared := self.RiichiState() == player.RiichiDeclared
	var candidates []actionCandidate
	for _, discard := range normalizedSelfTurnDiscards(actions) {
		afterDiscard, err := h.Discard(discard.Tile())
		if err != nil {
			return nil, fmt.Errorf("cannot build self-turn candidate for %s: %w", discard.Tile(), err)
		}
		shanten := candidateShanten(discard.Tile(), turnShanten, turnGoals)
		if riichi != nil && shanten <= 0 {
			// Match Manue's riichi candidate filtering: only regular-hand shanten
			// from AnalyzeShanten is considered here, so Seven Pairs and
			// Thirteen Orphans do not create riichi candidates.
			candidates = append(candidates, buildSelfTurnCandidate(riichi, discard.Tile(), h, afterDiscard, turnShanten, shanten, turnGoals, true, true))
		}
		if riichiDeclared && shanten > 0 {
			continue
		}
		candidates = append(candidates, buildSelfTurnCandidate(discard, discard.Tile(), h, afterDiscard, turnShanten, shanten, turnGoals, false, riichiDeclared))
	}
	return candidates, nil
}

func getOtherDiscardReactionCandidates(actions []action.Action, self player.PlayerViewer) ([]actionCandidate, error) {
	h, err := selfTurnHand(self)
	if err != nil {
		return nil, fmt.Errorf("cannot build reaction candidates: %w", err)
	}

	candidates := make([]actionCandidate, 0, len(actions))
	if pass := firstActionOfType[*action.Pass](actions); pass != nil {
		shanten, goals := service.AnalyzeShanten(h, service.AllowedExtraTiles(1))
		unknown := tile.MustTileFromCode("?")
		candidates = append(candidates, actionCandidate{
			traceKey:         "none",
			action:           pass,
			discardTile:      unknown,
			melds:            self.Melds(),
			turnHand:         h,
			afterDiscardHand: h,
			baseShanten:      shanten,
			shantenGoals:     goals,
			score:            scoreDiscardCandidate(unknown, shanten),
		})
	}

	callIndex := 0
	for _, a := range actions {
		callMeld, ok, err := actionCallMeld(a)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		callCandidates, err := getCallReactionCandidates(callIndex, a, callMeld, h, self.Melds())
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, callCandidates...)
		callIndex++
	}
	return candidates, nil
}

func getCallReactionCandidates(
	callIndex int,
	callAction action.Action,
	callMeld meld.Meld,
	baseHand *hand.VisibleHand,
	baseMelds []meld.Meld,
) ([]actionCandidate, error) {
	turnHand, err := baseHand.Call(callMeld)
	if err != nil {
		return nil, fmt.Errorf("cannot build reaction candidates for call %d: %w", callIndex, err)
	}
	nextMelds := append(slices.Clone(baseMelds), callMeld)
	turnShanten, turnGoals := service.AnalyzeShanten(turnHand, service.AllowedExtraTiles(1))

	if _, ok := callMeld.(*meld.CalledKan); ok {
		unknown := tile.MustTileFromCode("?")
		shanten, goals := service.AnalyzeShanten(turnHand, service.AllowedExtraTiles(1))
		return []actionCandidate{{
			traceKey:         fmt.Sprintf("%d.none", callIndex),
			action:           callAction,
			discardTile:      unknown,
			melds:            nextMelds,
			turnHand:         turnHand,
			afterDiscardHand: turnHand,
			baseShanten:      shanten,
			shantenGoals:     goals,
			score:            scoreDiscardCandidate(unknown, shanten),
		}}, nil
	}

	swapCallTiles := callSwapTiles(callMeld)
	discardTiles := tile.Tiles(turnHand.ToTiles()).Distinct(func(t tile.Tile) bool {
		return isSwapCallTile(t, swapCallTiles)
	})
	candidates := make([]actionCandidate, 0, len(discardTiles))
	for _, discardTile := range discardTiles {
		afterDiscard, err := turnHand.Discard(discardTile)
		if err != nil {
			return nil, fmt.Errorf("cannot build reaction candidate %d.%s: %w", callIndex, discardTile, err)
		}
		shanten := candidateShanten(discardTile, turnShanten, turnGoals)
		candidates = append(candidates, actionCandidate{
			traceKey:         fmt.Sprintf("%d.%s", callIndex, discardTile),
			action:           callAction,
			discardTile:      discardTile,
			melds:            nextMelds,
			turnHand:         turnHand,
			afterDiscardHand: afterDiscard,
			baseShanten:      turnShanten,
			shantenGoals:     turnGoals,
			score:            scoreDiscardCandidate(discardTile, shanten),
		})
	}
	return candidates, nil
}

func actionCallMeld(a action.Action) (meld.Meld, bool, error) {
	switch call := a.(type) {
	case *action.Chii:
		m, err := meld.NewChii(call.Taken(), call.Consumed(), call.Target())
		if err != nil {
			return nil, false, err
		}
		return m, true, nil
	case *action.Pon:
		m, err := meld.NewPon(call.Taken(), call.Consumed(), call.Target())
		if err != nil {
			return nil, false, err
		}
		return m, true, nil
	case *action.CalledKan:
		m, err := meld.NewCalledKan(call.Taken(), call.Consumed(), call.Target())
		if err != nil {
			return nil, false, err
		}
		return m, true, nil
	default:
		return nil, false, nil
	}
}

func callSwapTiles(m meld.Meld) []tile.Tile {
	switch c := m.(type) {
	case *meld.Chii:
		return c.SwapCallTiles()
	case *meld.Pon:
		return c.SwapCallTiles()
	default:
		return nil
	}
}

func isSwapCallTile(t tile.Tile, swapCallTiles []tile.Tile) bool {
	return slices.ContainsFunc(swapCallTiles, func(s tile.Tile) bool {
		return t.HasSameSymbol(s)
	})
}

// normalizedSelfTurnDiscards preserves CoffeeScript output behavior: when the
// same exact tile can be discarded from hand or as tsumogiri, keep tsumogiri.
func normalizedSelfTurnDiscards(actions []action.Action) []*action.Discard {
	discards := make([]*action.Discard, 0, len(actions))
	indexByTile := make(map[tile.Tile]int, len(actions))
	for _, a := range actions {
		discard, ok := a.(*action.Discard)
		if !ok {
			continue
		}
		if i, ok := indexByTile[discard.Tile()]; ok {
			if discard.Tsumogiri() {
				discards[i] = discard
			}
			continue
		}
		indexByTile[discard.Tile()] = len(discards)
		discards = append(discards, discard)
	}
	return discards
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
	baseShanten int,
	shanten int,
	goals []service.Goal,
	riichi bool,
	scoreAsRiichi bool,
) actionCandidate {
	return actionCandidate{
		traceKey:         formatDiscardTraceKey(riichi, discardTile),
		action:           immediateAction,
		riichi:           riichi,
		scoreAsRiichi:    scoreAsRiichi,
		discardTile:      discardTile,
		turnHand:         turnHand,
		afterDiscardHand: afterDiscardHand,
		baseShanten:      baseShanten,
		shantenGoals:     goals,
		score:            scoreDiscardCandidate(discardTile, shanten),
	}
}

func scoreDiscardCandidate(discardTile tile.Tile, shanten int) candidateScore {
	return candidateScore{
		avgRank: 0,
		expPts:  0,
		shanten: shanten,
		red:     discardTile.IsRed(),
	}
}
