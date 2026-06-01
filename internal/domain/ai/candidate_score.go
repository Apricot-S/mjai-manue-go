package ai

import (
	"fmt"
	"slices"

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
}

func chooseBestCandidate(candidates []actionCandidate, preferBlack bool) actionCandidate {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if compareCandidates(candidate, best, preferBlack) < 0 {
			best = candidate
		}
	}
	return best
}

func sortedCandidates(candidates []actionCandidate, preferBlack bool) []actionCandidate {
	sortedCandidates := slices.Clone(candidates)
	slices.SortFunc(sortedCandidates, func(lhs, rhs actionCandidate) int {
		return compareCandidates(lhs, rhs, preferBlack)
	})
	return sortedCandidates
}

func compareCandidates(lhs, rhs actionCandidate, preferBlack bool) int {
	if result := compareCandidateScore(&lhs.score, &rhs.score); result != 0 {
		return result
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

func compareCandidateScore(lhs, rhs *candidateScore) int {
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
	return 0
}

func evaluateCandidateScore(
	score candidateScore,
	scoreChanges scoreDeltaProbDist,
	stats RankStats,
	state rankStateViewer,
	self seat.Seat,
) candidateScore {
	scores := state.Scores()
	startingDealer := state.StartingDealer()
	score.expectedPoints = scoreChanges.expected()[self.Index()]
	score.averageRank = averageRank(
		scoreChanges,
		self.Index(),
		float64(scores[self.Index()]),
		self.DistanceFrom(startingDealer),
		buildRankOpponents(stats, state, self),
	)
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
		score.exhaustiveDrawProb,
		otherWinDists,
		score.otherWinProb,
	)
	return immediateDist.replace(scoreDelta{}, futureDist)
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
	return evaluateCandidateScore(score, scoreChanges, rankStats, state, self), nil
}
