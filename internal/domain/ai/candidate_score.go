package ai

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

type candidateScore struct {
	// averageRank is the average rank.
	averageRank float64
	// expectedPoints is the expected points.
	expectedPoints float64
	// dealInProb is the deal-in probability.
	dealInProb float64
	// winProb is the win probability.
	winProb float64
	// exhaustiveDrawProb is the exhaustive draw probability.
	exhaustiveDrawProb float64
	// otherWinProb is the other players' win probability.
	otherWinProb float64
	// averageWinPoints is the average win points.
	averageWinPoints float64
	// exhaustiveDrawAveragePoints is the average exhaustive draw points.
	exhaustiveDrawAveragePoints float64
	// shanten is the shanten number.
	shanten int
	// red indicates whether the candidate discards a red tile.
	red bool
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
	if lhs.averageRank < rhs.averageRank {
		return -1
	}
	if lhs.averageRank > rhs.averageRank {
		return 1
	}
	if lhs.expectedPoints > rhs.expectedPoints {
		return -1
	}
	if lhs.expectedPoints < rhs.expectedPoints {
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
	score.expectedPoints = expectedPts(selfID, scoreChanges)
	score.averageRank = averageRank(scoreChanges, selfID, selfScore, selfPosition, opponents)
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
		score.exhaustiveDrawProb,
		otherWinDists,
		score.otherWinProb,
	)
	return totalScoreDeltaDist(immediateDist, futureDist)
}

func evaluateCandidateFromComponents(
	score candidateScore,
	dealInEstimates []dealInEstimate,
	winEstimate winEstimate,
	exhaustiveDrawProbOnSelfNoWin float64,
	exhaustiveDrawAveragePoints float64,
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
	score.averageWinPoints = winEstimate.avgPts
	exhaustiveDrawProb, otherWinProb, err := remainingRoundEndProbs(score.winProb, exhaustiveDrawProbOnSelfNoWin)
	if err != nil {
		return candidateScore{}, err
	}
	score.exhaustiveDrawProb = exhaustiveDrawProb
	score.otherWinProb = otherWinProb
	score.exhaustiveDrawAveragePoints = exhaustiveDrawAveragePoints
	scoreChanges := candidateTotalScoreDeltaDist(
		score,
		immediateDist,
		selfWinDist,
		exhaustiveDrawDist,
		otherWinDists,
	)
	return evaluateCandidateScoreFromState(score, scoreChanges, rankStats, state, self), nil
}
