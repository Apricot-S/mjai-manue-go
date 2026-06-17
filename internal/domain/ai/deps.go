package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

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
	DealInStats
	RankStats
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

type DrawTenpaiStats interface {
	ExhaustiveDrawNotenCount() int
	ExhaustiveDrawTenpaiTurnFreq(turnKey string) (freq int, ok bool)
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

type DangerEstimator interface {
	EstimateDealInProb(state round.StateViewer, self seat.Seat, winner seat.Seat, discard tile.Tile) (float64, error)
}
