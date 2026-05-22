package ai

import (
	"fmt"
	"maps"
	"math"
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
)

type dealInEstimate struct {
	winnerID int
	prob     float64
}

type winEstimate struct {
	// prob is the win probability.
	prob float64
	// avgPts is the average win points when self wins.
	avgPts float64
	// expPts is the expected win points over all estimation trials.
	expPts float64
	// pointsDist is the win-points distribution when self wins.
	pointsDist scalarProbDist
}

type winEstimateAccumulator struct {
	numTries   int
	totalWins  int
	totalPts   float64
	pointFreqs map[float64]int
}

func (a *winEstimateAccumulator) addNoWinTrial() {
	a.numTries++
}

func (a *winEstimateAccumulator) addWinTrial(points float64) error {
	if points <= 0 {
		return fmt.Errorf("cannot add win trial: points must be positive")
	}
	a.numTries++
	a.totalWins++
	a.totalPts += points
	if a.pointFreqs == nil {
		a.pointFreqs = make(map[float64]int)
	}
	a.pointFreqs[points]++
	return nil
}

func (a *winEstimateAccumulator) merge(other winEstimateAccumulator) error {
	if other.numTries < 0 || other.totalWins < 0 || other.totalPts < 0 {
		return fmt.Errorf("cannot merge win estimate accumulator: values must be non-negative")
	}
	if other.totalWins > other.numTries {
		return fmt.Errorf("cannot merge win estimate accumulator: totalWins must not exceed numTries")
	}

	countedWins := 0
	for points, freq := range other.pointFreqs {
		if points <= 0 {
			return fmt.Errorf("cannot merge win estimate accumulator: points must be positive")
		}
		if freq < 0 {
			return fmt.Errorf("cannot merge win estimate accumulator: point frequency must be non-negative")
		}
		countedWins += freq
	}
	if countedWins != other.totalWins {
		return fmt.Errorf("cannot merge win estimate accumulator: point frequencies must sum to totalWins")
	}

	a.numTries += other.numTries
	a.totalWins += other.totalWins
	a.totalPts += other.totalPts
	if len(other.pointFreqs) > 0 && a.pointFreqs == nil {
		a.pointFreqs = make(map[float64]int, len(other.pointFreqs))
	}
	for points, freq := range other.pointFreqs {
		a.pointFreqs[points] += freq
	}
	return nil
}

func (a winEstimateAccumulator) estimate() (winEstimate, error) {
	return newWinEstimate(a.numTries, a.totalWins, a.totalPts, a.pointFreqs)
}

func (a winEstimateAccumulator) clone() winEstimateAccumulator {
	return winEstimateAccumulator{
		numTries:   a.numTries,
		totalWins:  a.totalWins,
		totalPts:   a.totalPts,
		pointFreqs: maps.Clone(a.pointFreqs),
	}
}

type winEstimateAccumulatorSet map[string]winEstimateAccumulator

func newWinEstimateAccumulatorSet(keys []string) winEstimateAccumulatorSet {
	accumulators := make(winEstimateAccumulatorSet, len(keys))
	for _, key := range keys {
		accumulators[key] = winEstimateAccumulator{}
	}
	return accumulators
}

func (s winEstimateAccumulatorSet) addNoWinTrial(key string) {
	accumulator := s[key]
	accumulator.addNoWinTrial()
	s[key] = accumulator
}

func (s winEstimateAccumulatorSet) addWinTrial(key string, points float64) error {
	accumulator := s[key]
	if err := accumulator.addWinTrial(points); err != nil {
		return err
	}
	s[key] = accumulator
	return nil
}

func (s winEstimateAccumulatorSet) addTrial(winPtsByKey map[string]float64) error {
	for key, points := range winPtsByKey {
		if _, ok := s[key]; !ok {
			return fmt.Errorf("cannot add win estimate trial: unknown candidate key %q", key)
		}
		if points <= 0 {
			return fmt.Errorf("cannot add win estimate trial: points must be positive")
		}
	}

	for key := range s {
		points, ok := winPtsByKey[key]
		if ok {
			if err := s.addWinTrial(key, points); err != nil {
				return err
			}
			continue
		}
		s.addNoWinTrial(key)
	}
	return nil
}

func (s winEstimateAccumulatorSet) merge(other winEstimateAccumulatorSet) error {
	merged := make(winEstimateAccumulatorSet, len(other))
	for key, otherAccumulator := range other {
		accumulator := s[key].clone()
		if err := accumulator.merge(otherAccumulator); err != nil {
			return fmt.Errorf("cannot merge win estimate accumulator for %q: %w", key, err)
		}
		merged[key] = accumulator
	}
	maps.Copy(s, merged)
	return nil
}

func (s winEstimateAccumulatorSet) estimates() (map[string]winEstimate, error) {
	estimates := make(map[string]winEstimate, len(s))
	for key, accumulator := range s {
		estimate, err := accumulator.estimate()
		if err != nil {
			return nil, fmt.Errorf("cannot build win estimate for %q: %w", key, err)
		}
		estimates[key] = estimate
	}
	return estimates, nil
}

func winEstimatesFromTrials(keys []string, trials []map[string]float64) (map[string]winEstimate, error) {
	accumulators := newWinEstimateAccumulatorSet(keys)
	for _, trial := range trials {
		if err := accumulators.addTrial(trial); err != nil {
			return nil, err
		}
	}
	return accumulators.estimates()
}

func newWinEstimate(
	numTries int,
	totalWins int,
	totalPts float64,
	pointFreqs map[float64]int,
) (winEstimate, error) {
	if numTries <= 0 {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: numTries must be positive")
	}
	if totalWins < 0 {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: totalWins must be non-negative")
	}
	if totalWins > numTries {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: totalWins must not exceed numTries")
	}
	if totalPts < 0 {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: totalPts must be non-negative")
	}
	if totalWins == 0 {
		if totalPts != 0 || len(pointFreqs) != 0 {
			return winEstimate{}, fmt.Errorf("cannot build win estimate: zero wins must have no points")
		}
		return winEstimate{
			prob:       0,
			avgPts:     0,
			expPts:     0,
			pointsDist: scalarProbDist{},
		}, nil
	}

	totalWinsFloat := float64(totalWins)
	dist := make(map[float64]float64, len(pointFreqs))
	countedWins := 0
	for points, freq := range pointFreqs {
		if points <= 0 {
			return winEstimate{}, fmt.Errorf("cannot build win estimate: points must be positive")
		}
		if freq < 0 {
			return winEstimate{}, fmt.Errorf("cannot build win estimate: point frequency must be non-negative")
		}
		countedWins += freq
		dist[points] = float64(freq) / totalWinsFloat
	}
	if countedWins != totalWins {
		return winEstimate{}, fmt.Errorf("cannot build win estimate: point frequencies must sum to totalWins")
	}

	return winEstimate{
		prob:       totalWinsFloat / float64(numTries),
		avgPts:     totalPts / totalWinsFloat,
		expPts:     totalPts / float64(numTries),
		pointsDist: newScalarProbDist(dist),
	}, nil
}

func dealInProb(estimates []dealInEstimate) (float64, error) {
	safeProb, err := safeProb(estimates)
	if err != nil {
		return 0, err
	}
	return 1.0 - safeProb, nil
}

func safeProb(estimates []dealInEstimate) (float64, error) {
	safeProb := 1.0
	for _, estimate := range estimates {
		if estimate.prob < 0.0 || estimate.prob > 1.0 {
			return 0, fmt.Errorf("cannot estimate safe probability: deal-in probability must be between 0 and 1")
		}
		safeProb *= 1.0 - estimate.prob
	}
	return safeProb, nil
}

// winScoreFactor returns how one win point unit changes all players' scores.
//
// actorID is the winner. targetID is the winner for self-draw wins, or the
// discarder for ron wins. dealerID is the round dealer.
func winScoreFactor(actorID int, targetID int, dealerID int) scoreDelta {
	if targetID != actorID {
		// Ron: the discarder pays the full win points.
		var factor scoreDelta
		factor[actorID] = 1.0
		factor[targetID] = -1.0
		return factor
	}

	if actorID == dealerID {
		// Dealer self-draw: each non-dealer pays one third.
		factor := scoreDelta{-1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0, -1.0 / 3.0}
		factor[actorID] = 1.0
		return factor
	}

	// Non-dealer self-draw: the dealer pays half, each other non-dealer pays a quarter.
	factor := scoreDelta{-1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0, -1.0 / 4.0}
	factor[actorID] = 1.0
	factor[dealerID] = -1.0 / 2.0
	return factor
}

// winScoreFactorDist returns the distribution of score factors for a winner.
//
// selfDrawProb is the probability that the win is by self draw. Ron targets are
// treated as uniformly distributed among the three other players.
func winScoreFactorDist(actorID int, dealerID int, selfDrawProb float64) scoreDeltaProbDist {
	dist := make(scoreDeltaProbDist, 4)
	ronTargetProb := (1.0 - selfDrawProb) / 3.0
	for targetID := range 4 {
		var prob float64
		if targetID == actorID {
			prob = selfDrawProb
		} else {
			prob = ronTargetProb
		}
		dist[winScoreFactor(actorID, targetID, dealerID)] = prob
	}
	return newScoreDeltaProbDist(dist)
}

// winPointsDist returns a probability distribution from win-points frequencies.
func winPointsDist(pointFreqs map[string]int) (scalarProbDist, error) {
	totalFreqs := pointFreqs["total"]
	if totalFreqs <= 0 {
		return nil, fmt.Errorf("invalid win points frequencies: total must be positive")
	}
	totalFreqsFloat := float64(totalFreqs)

	dist := make(map[float64]float64, len(pointFreqs)-1)
	for points, freq := range pointFreqs {
		if points == "total" {
			continue
		}
		parsedPoints, err := strconv.ParseFloat(points, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid win points frequency key %q: %w", points, err)
		}
		dist[parsedPoints] = float64(freq) / totalFreqsFloat
	}
	return newScalarProbDist(dist), nil
}

// randomWinScoreDeltaDist returns the score-change distribution for a random
// win by actorID.
func randomWinScoreDeltaDist(
	actorID int,
	dealerID int,
	selfDrawProb float64,
	pointFreqs map[string]int,
) (scoreDeltaProbDist, error) {
	pointsDist, err := winPointsDist(pointFreqs)
	if err != nil {
		return nil, err
	}
	return multiplyScalarScoreDeltaProbDists(
		pointsDist,
		winScoreFactorDist(actorID, dealerID, selfDrawProb),
	), nil
}

func randomWinScoreDeltaDistFromStats(
	actorID int,
	dealerID int,
	stats WinScoreStats,
) (scoreDeltaProbDist, error) {
	if stats.NumWins() <= 0 {
		return nil, fmt.Errorf("cannot build random win score delta distribution: numWins must be positive")
	}

	pointFreqs := stats.NonDealerWinPointFreqs()
	if actorID == dealerID {
		pointFreqs = stats.DealerWinPointFreqs()
	}
	return randomWinScoreDeltaDist(
		actorID,
		dealerID,
		float64(stats.NumSelfDrawWins())/float64(stats.NumWins()),
		pointFreqs,
	)
}

func winScoreDeltaDistFromPointsDist(
	actorID int,
	dealerID int,
	stats WinScoreStats,
	pointsDist scalarProbDist,
) (scoreDeltaProbDist, error) {
	if stats.NumWins() <= 0 {
		return nil, fmt.Errorf("cannot build win score delta distribution: numWins must be positive")
	}
	return multiplyScalarScoreDeltaProbDists(
		pointsDist,
		winScoreFactorDist(actorID, dealerID, float64(stats.NumSelfDrawWins())/float64(stats.NumWins())),
	), nil
}

func selfWinScoreDeltaDistFromEstimate(
	selfID int,
	dealerID int,
	stats WinScoreStats,
	estimate winEstimate,
) (scoreDeltaProbDist, error) {
	return winScoreDeltaDistFromPointsDist(selfID, dealerID, stats, estimate.pointsDist)
}

func exhaustiveDrawProb(stats RoundEndStats, currentTurn float64) (float64, error) {
	currentTurnIndex := int(currentTurn)
	turnDistribution := stats.TurnDistribution()
	if currentTurnIndex < 0 || currentTurnIndex >= numTurnDistributionEntries {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw probability: current turn is out of range")
	}

	remainingProb := 0.0
	for _, prob := range turnDistribution[currentTurnIndex:] {
		remainingProb += prob
	}
	if remainingProb <= 0 {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw probability: remaining turn probability must be positive")
	}
	return stats.ExhaustiveDrawRatio() / remainingProb, nil
}

func exhaustiveDrawProbOnSelfNoWin(stats RoundEndStats, currentTurn float64) (float64, error) {
	prob, err := exhaustiveDrawProb(stats, currentTurn)
	if err != nil {
		return 0, err
	}
	return math.Pow(prob, 3.0/4.0), nil
}

func expectedRemainingTurns(stats RoundEndStats, currentTurn float64) (int, error) {
	currentTurnIndex := int(math.Round(currentTurn))
	if currentTurnIndex < 0 || currentTurnIndex >= numTurnDistributionEntries {
		return 0, fmt.Errorf("cannot estimate expected remaining turns: current turn is out of range")
	}

	turnDistribution := stats.TurnDistribution()
	num := 0.0
	den := 0.0
	for i := currentTurnIndex; i < numTurnDistributionEntries; i++ {
		prob := turnDistribution[i]
		num += prob * (float64(i) - math.Round(currentTurn) + 0.5)
		den += prob
	}
	if den == 0 {
		return 0, nil
	}
	return int(math.Round(num / den)), nil
}

func tenpaiProb(stats TenpaiEstimatorStats, riichi bool, remainTurns int, numMelds int) (float64, error) {
	if riichi {
		return 1.0, nil
	}

	total, tenpai, ok := stats.YamitenCounts(remainTurns, numMelds)
	if !ok {
		return 1.0, nil
	}
	if total <= 0 {
		return 0, fmt.Errorf("cannot estimate tenpai probability: yamiten total must be positive")
	}
	if tenpai < 0 || tenpai > total {
		return 0, fmt.Errorf("cannot estimate tenpai probability: yamiten tenpai must be between 0 and total")
	}
	return float64(tenpai) / float64(total), nil
}

// dealInExpPts returns expected points from dealing in. It is negative because
// it represents self's score change when the discard is unsafe.
func dealInExpPts(stats DealInStats, safeProb float64) (float64, error) {
	if safeProb < 0.0 || safeProb > 1.0 {
		return 0.0, fmt.Errorf("cannot estimate deal-in expected points: safe probability must be between 0 and 1")
	}
	return -(1.0 - safeProb) * stats.AvgWinPts(), nil
}

// immediateDealInScoreDeltaDist returns the immediate score-change
// distribution for a discard against one possible ron winner.
//
// dealInProb is the probability of dealing in to winnerID. The no-change branch
// means this opponent did not ron the discard and the round can continue.
func immediateDealInScoreDeltaDist(
	winnerID int,
	selfID int,
	dealInProb float64,
	pointsDist scalarProbDist,
) (scoreDeltaProbDist, error) {
	if dealInProb < 0.0 || dealInProb > 1.0 {
		return nil, fmt.Errorf("cannot build immediate deal-in score delta distribution: deal-in probability must be between 0 and 1")
	}

	var dealInFactor scoreDelta
	dealInFactor[winnerID] = 1.0
	dealInFactor[selfID] = -1.0
	unitDist := newScoreDeltaProbDist(map[scoreDelta]float64{
		dealInFactor: dealInProb,
		{}:           1.0 - dealInProb,
	})
	return multiplyScalarScoreDeltaProbDists(pointsDist, unitDist), nil
}

func immediateDealInScoreDeltaDistFromStats(
	winnerID int,
	selfID int,
	dealerID int,
	dealInProb float64,
	stats WinScoreStats,
) (scoreDeltaProbDist, error) {
	pointFreqs := stats.NonDealerWinPointFreqs()
	if winnerID == dealerID {
		pointFreqs = stats.DealerWinPointFreqs()
	}
	pointsDist, err := winPointsDist(pointFreqs)
	if err != nil {
		return nil, err
	}
	return immediateDealInScoreDeltaDist(winnerID, selfID, dealInProb, pointsDist)
}

func immediateScoreDeltaDistFromStats(
	selfID int,
	dealerID int,
	estimates []dealInEstimate,
	stats WinScoreStats,
) (scoreDeltaProbDist, error) {
	dists := make([]scoreDeltaProbDist, 0, len(estimates))
	for _, estimate := range estimates {
		dist, err := immediateDealInScoreDeltaDistFromStats(
			estimate.winnerID,
			selfID,
			dealerID,
			estimate.prob,
			stats,
		)
		if err != nil {
			return nil, err
		}
		dists = append(dists, dist)
	}
	return immediateScoreDeltaDist(dists), nil
}

// immediateScoreDeltaDist merges immediate deal-in distributions from multiple
// possible ron winners.
//
// Each distribution must include a no-change branch. The merge replaces only
// that branch, so it models the same first-ron approximation as Manue and avoids
// expanding double/triple ron combinations.
func immediateScoreDeltaDist(dists []scoreDeltaProbDist) scoreDeltaProbDist {
	result := scoreDeltaProbDist{{}: 1.0}
	for _, dist := range dists {
		result = result.replace(scoreDelta{}, dist)
	}
	return result
}

func safeWinExpPts(safeProb float64, avgWinPts float64) (float64, error) {
	if safeProb < 0.0 || safeProb > 1.0 {
		return 0.0, fmt.Errorf("cannot estimate safe win expected points: safe probability must be between 0 and 1")
	}
	return safeProb * avgWinPts, nil
}

func exhaustiveDrawExpPts(safeProb float64, exhaustiveDrawProb float64, avgDrawPts float64) (float64, error) {
	if safeProb < 0.0 || safeProb > 1.0 {
		return 0.0, fmt.Errorf("cannot estimate exhaustive-draw expected points: safe probability must be between 0 and 1")
	}
	if exhaustiveDrawProb < 0.0 || exhaustiveDrawProb > 1.0 {
		return 0.0, fmt.Errorf("cannot estimate exhaustive-draw expected points: exhaustive-draw probability must be between 0 and 1")
	}
	return safeProb * exhaustiveDrawProb * avgDrawPts, nil
}

// remainingRoundEndProbs splits the probability mass left after accounting for
// self's future win into exhaustive draw and other players' win.
func remainingRoundEndProbs(selfWinProb float64, exhaustiveDrawProbOnSelfNoWin float64) (drawProb float64, othersWinProb float64, err error) {
	if selfWinProb < 0.0 || selfWinProb > 1.0 {
		return 0.0, 0.0, fmt.Errorf("cannot estimate remaining round-end probabilities: self win probability must be between 0 and 1")
	}
	if exhaustiveDrawProbOnSelfNoWin < 0.0 || exhaustiveDrawProbOnSelfNoWin > 1.0 {
		return 0.0, 0.0, fmt.Errorf("cannot estimate remaining round-end probabilities: exhaustive-draw probability must be between 0 and 1")
	}

	noSelfWinProb := 1.0 - selfWinProb
	return noSelfWinProb * exhaustiveDrawProbOnSelfNoWin,
		noSelfWinProb * (1.0 - exhaustiveDrawProbOnSelfNoWin),
		nil
}

// futureScoreDeltaDist mixes possible future round-ending outcomes after a
// discard did not immediately deal in.
func futureScoreDeltaDist(
	selfWinDist scoreDeltaProbDist,
	selfWinProb float64,
	exhaustiveDrawDist scoreDeltaProbDist,
	exhaustiveDrawProb float64,
	otherWinDists []scoreDeltaProbDist,
	othersWinProb float64,
) scoreDeltaProbDist {
	items := []weightedScoreDeltaProbDist{
		{dist: selfWinDist, prob: selfWinProb},
		{dist: exhaustiveDrawDist, prob: exhaustiveDrawProb},
	}
	otherWinProb := 0.0
	if len(otherWinDists) > 0 {
		otherWinProb = othersWinProb / float64(len(otherWinDists))
	}
	for _, dist := range otherWinDists {
		items = append(items, weightedScoreDeltaProbDist{dist: dist, prob: otherWinProb})
	}
	return mergeScoreDeltaProbDists(items)
}

// totalScoreDeltaDist replaces the no-change branch of immediateDist with the
// future round-ending distribution. This mirrors Manue's flow where no
// immediate deal-in means the round continues.
func totalScoreDeltaDist(immediateDist scoreDeltaProbDist, futureDist scoreDeltaProbDist) scoreDeltaProbDist {
	return immediateDist.replace(scoreDelta{}, futureDist)
}

func expectedPts(selfID int, scoreChanges scoreDeltaProbDist) float64 {
	return scoreChanges.expected()[selfID]
}

// notenExhaustiveDrawTenpaiProb returns the probability that a currently
// noten player reaches tenpai before exhaustive draw, conditional on the round
// ending by exhaustive draw.
func notenExhaustiveDrawTenpaiProb(stats DrawTenpaiStats, currentTurn float64) (float64, error) {
	notenFreq := stats.ExhaustiveDrawNotenCount()
	if notenFreq < 0 {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw tenpai probability: noten count must be non-negative")
	}

	tenpaiFreq := 0
	for turn := currentTurn + 0.25; turn <= round.FinalTurn; turn += 0.25 {
		key := strconv.FormatFloat(turn, 'f', -1, 64)
		freq, ok := stats.ExhaustiveDrawTenpaiTurnFreq(key)
		if !ok {
			return 0, fmt.Errorf("cannot estimate exhaustive-draw tenpai probability: missing tenpai turn frequency for turn %s", key)
		}
		tenpaiFreq += freq
	}

	totalFreq := tenpaiFreq + notenFreq
	if totalFreq <= 0 {
		return 0, fmt.Errorf("cannot estimate exhaustive-draw tenpai probability: frequency total must be positive")
	}
	return float64(tenpaiFreq) / float64(totalFreq), nil
}

// exhaustiveDrawTenpaiProbs adjusts current tenpai probabilities into
// exhaustive-draw tenpai probabilities. It keeps the current tenpai probability
// and adds the chance that a currently noten player reaches tenpai before
// exhaustive draw.
func exhaustiveDrawTenpaiProbs(currentTenpaiProbs [4]float64, notenTenpaiProb float64) [4]float64 {
	var probs [4]float64
	for playerID, currentTenpaiProb := range currentTenpaiProbs {
		probs[playerID] = currentTenpaiProb + (1.0-currentTenpaiProb)*notenTenpaiProb
	}
	return probs
}

// ryukyokuScoreDelta returns the score change vector for exhaustive draw
// tenpai payments.
func ryukyokuScoreDelta(tenpais [4]bool) scoreDelta {
	points := service.RyukyokuPoints(tenpais)
	var delta scoreDelta
	for i, point := range points {
		delta[i] = float64(point)
	}
	return delta
}

// exhaustiveDrawScoreDeltaDistFromTenpaiProbs returns the score change
// distribution assuming the round ends in an exhaustive draw.
func exhaustiveDrawScoreDeltaDistFromTenpaiProbs(tenpaiProbs [4]float64) scoreDeltaProbDist {
	tenpaisDist := aheadVectorProbDist{{}: 1.0}
	for playerID, tenpaiProb := range tenpaiProbs {
		var tenpais aheadVector
		tenpais[playerID] = 1
		tenpaisDist = addAheadVectorProbDists(tenpaisDist, newAheadVectorProbDist(map[aheadVector]float64{
			{}:      1.0 - tenpaiProb,
			tenpais: tenpaiProb,
		}))
	}

	return tenpaisDist.mapValueScoreDelta(func(tenpais aheadVector) scoreDelta {
		return ryukyokuScoreDelta(aheadVectorToBoolArray(tenpais))
	})
}

// futureExhaustiveDrawScoreDeltaDist returns the exhaustive-draw score change
// distribution from current tenpai probabilities. It first adjusts current
// tenpai probabilities by the chance that currently noten players reach tenpai
// before exhaustive draw.
func futureExhaustiveDrawScoreDeltaDist(
	currentTenpaiProbs [4]float64,
	notenTenpaiProb float64,
) scoreDeltaProbDist {
	return exhaustiveDrawScoreDeltaDistFromTenpaiProbs(exhaustiveDrawTenpaiProbs(currentTenpaiProbs, notenTenpaiProb))
}

// exhaustiveDrawAvgPts returns self's expected score change assuming the round
// ends in an exhaustive draw and tenpaiProbs already represent exhaustive-draw
// tenpai probabilities.
func exhaustiveDrawAvgPts(selfID int, tenpaiProbs [4]float64) float64 {
	return exhaustiveDrawScoreDeltaDistFromTenpaiProbs(tenpaiProbs).expected()[selfID]
}

func aheadVectorToBoolArray(value aheadVector) [4]bool {
	var result [4]bool
	for i, v := range value {
		result[i] = v != 0
	}
	return result
}
