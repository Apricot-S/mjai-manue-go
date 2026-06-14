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

func evaluateCandidateFromComponents(
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
	var score candidateScore
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
	futureDist := futureScoreDeltaDist(
		selfWinDist,
		score.winProb,
		exhaustiveDrawDist,
		score.exhaustiveDrawProb,
		otherWinDists,
		score.otherWinProb,
	)
	scoreChanges := immediateDist.replace(scoreDelta{}, futureDist)
	scores := state.Scores()
	startingDealer := state.StartingDealer()
	score.expectedPoints = scoreChanges.expected()[self.Index()]
	score.averageRank = averageRank(
		scoreChanges,
		self.Index(),
		float64(scores[self.Index()]),
		self.DistanceFrom(startingDealer),
		buildRankOpponents(rankStats, state, self),
	)
	return score, nil
}
