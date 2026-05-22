package ai

import (
	"strconv"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type relativeWinProbTable map[string]float64

type rankStateViewer interface {
	NextRound() (wind.Wind, int)
	Scores() [common.NumPlayers]int
	StartingDealer() seat.Seat
}

type rankOpponent struct {
	id       int
	score    float64
	position int
	winProbs relativeWinProbTable
}

func relativeWinProbs(
	stats RankStats,
	roundWind wind.Wind,
	roundNumber int,
	selfPosition int,
	otherPosition int,
) relativeWinProbTable {
	winProbs, ok := stats.RelativeWinProbs(roundWind, roundNumber, selfPosition, otherPosition)
	if !ok {
		return nil
	}
	return relativeWinProbTable(winProbs)
}

func buildRankOpponents(stats RankStats, state rankStateViewer, self seat.Seat) []rankOpponent {
	nextRoundWind, nextRoundNumber := state.NextRound()
	scores := state.Scores()
	startingDealer := state.StartingDealer()
	selfPosition := self.DistanceFrom(startingDealer)

	opponents := make([]rankOpponent, 0, common.NumPlayers-1)
	for i := range common.NumPlayers {
		opponentSeat := seat.MustSeat(i)
		if opponentSeat == self {
			continue
		}
		opponentPosition := opponentSeat.DistanceFrom(startingDealer)
		opponents = append(opponents, rankOpponent{
			id:       i,
			score:    float64(scores[i]),
			position: opponentPosition,
			winProbs: relativeWinProbs(stats, nextRoundWind, nextRoundNumber, selfPosition, opponentPosition),
		})
	}
	return opponents
}

// averageRank returns self's expected final rank from pairwise win probabilities
// against the other players.
func averageRank(
	scoreChanges scoreDeltaProbDist,
	selfID int,
	selfScore float64,
	selfPosition int,
	opponents []rankOpponent,
) float64 {
	winsDist := aheadVectorProbDist{{}: 1.0}
	for _, opponent := range opponents {
		winProb := winProbAgainst(
			scoreChanges,
			selfID,
			opponent.id,
			selfScore,
			opponent.score,
			selfPosition,
			opponent.position,
			opponent.winProbs,
		)

		var wins aheadVector
		wins[opponent.id] = 1
		winsDist = addAheadVectorProbDists(winsDist, newAheadVectorProbDist(map[aheadVector]float64{
			{}:   1.0 - winProb,
			wins: winProb,
		}))
	}

	rankDist := winsDist.mapValueScalar(func(wins aheadVector) float64 {
		return float64(4 - countAheadWins(wins))
	})
	return rankDist.expected()
}

// winProbAgainst returns the probability that self finishes ahead of another
// player after applying a score-delta distribution.
func winProbAgainst(
	scoreChanges scoreDeltaProbDist,
	selfID int,
	otherID int,
	selfScore float64,
	otherScore float64,
	selfPosition int,
	otherPosition int,
	winProbs relativeWinProbTable,
) float64 {
	relativeScoreDist := scoreChanges.mapValueScalar(func(scoreChange scoreDelta) float64 {
		return (selfScore + scoreChange[selfID]) - (otherScore + scoreChange[otherID])
	})

	winProb := 0.0
	for relativeScore, prob := range relativeScoreDist {
		winProb += prob * winProbFromRelativeScore(
			relativeScore,
			winProbs,
			selfPosition,
			otherPosition,
		)
	}
	return winProb
}

// winProbFromRelativeScore returns the probability that the player finishes
// ahead of another player from their relative score.
//
// When statistical data is missing, ties are broken by seating order from the
// starting dealer. The player closer to the starting dealer wins a tie.
func winProbFromRelativeScore(
	relativeScore float64,
	winProbs relativeWinProbTable,
	selfPosition int,
	otherPosition int,
) float64 {
	if winProbs != nil {
		if prob, ok := winProbs[strconv.FormatFloat(relativeScore, 'f', 0, 64)]; ok {
			return prob
		}
	}

	if selfPosition < otherPosition {
		if relativeScore >= 0.0 {
			return 1.0
		}
		return 0.0
	}
	if relativeScore > 0.0 {
		return 1.0
	}
	return 0.0
}
