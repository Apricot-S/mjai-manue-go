package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type StartRound struct {
	roundWind     wind.Wind
	roundNumber   int
	honba         int
	riichiDeposit int
	dealer        seat.Seat
	doraIndicator tile.Tile
	scores        *[common.NumPlayers]int
	hands         [common.NumPlayers][common.InitHandSize]tile.Tile
}

func NewStartRound(
	roundWind wind.Wind,
	roundNumber int,
	honba int,
	riichiDeposit int,
	dealer seat.Seat,
	doraIndicator tile.Tile,
	scores *[common.NumPlayers]int,
	hands [common.NumPlayers][common.InitHandSize]tile.Tile,
) *StartRound {
	return &StartRound{
		roundWind,
		roundNumber,
		honba,
		riichiDeposit,
		dealer,
		doraIndicator,
		scores,
		hands,
	}
}

func (*StartRound) isEvent() {}

func (s *StartRound) RoundWind() wind.Wind {
	return s.roundWind
}

func (s *StartRound) RoundNumber() int {
	return s.roundNumber
}

func (s *StartRound) Honba() int {
	return s.honba
}

func (s *StartRound) RiichiDeposit() int {
	return s.riichiDeposit
}

func (s *StartRound) Dealer() seat.Seat {
	return s.dealer
}

func (s *StartRound) DoraIndicator() tile.Tile {
	return s.doraIndicator
}

func (s *StartRound) Scores() *[common.NumPlayers]int {
	return s.scores
}

func (s *StartRound) Hands() [common.NumPlayers][common.InitHandSize]tile.Tile {
	return s.hands
}
