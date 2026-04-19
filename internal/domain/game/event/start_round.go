package event

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

const (
	minRoundNumber = 1
	maxRoundNumber = 4
	initHandSize   = 13
)

type StartRound struct {
	roundWind      wind.Wind
	roundNumber    int
	honba          int
	riichiDeposit  int
	dealer         seat.Seat
	startingDealer seat.Seat
	doraIndicator  tile.Tile
	scores         *[common.NumPlayers]int
	hands          [common.NumPlayers][initHandSize]tile.Tile
}

func NewStartRound(
	roundWind wind.Wind,
	roundNumber int,
	honba int,
	riichiDeposit int,
	dealer seat.Seat,
	startingDealer seat.Seat,
	doraIndicator tile.Tile,
	scores *[common.NumPlayers]int,
	hands [common.NumPlayers][initHandSize]tile.Tile,
) (*StartRound, error) {
	if roundWind < wind.East || wind.North < roundWind {
		return nil, fmt.Errorf("invalid round wind: %v", roundWind)
	}
	if roundNumber < minRoundNumber || maxRoundNumber < roundNumber {
		return nil, fmt.Errorf("invalid round number: %d", roundNumber)
	}
	if honba < 0 {
		return nil, fmt.Errorf("invalid honba: %d", honba)
	}
	if riichiDeposit < 0 {
		return nil, fmt.Errorf("invalid riichi deposit: %d", riichiDeposit)
	}
	if doraIndicator.IsUnknown() {
		return nil, fmt.Errorf("invalid dora indicator: %v", doraIndicator)
	}

	return &StartRound{
		roundWind,
		roundNumber,
		honba,
		riichiDeposit,
		dealer,
		startingDealer,
		doraIndicator,
		scores,
		hands,
	}, nil
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

func (s *StartRound) StartingDealer() seat.Seat {
	return s.startingDealer
}

func (s *StartRound) DoraIndicator() tile.Tile {
	return s.doraIndicator
}

func (s *StartRound) Scores() *[common.NumPlayers]int {
	return s.scores
}

func (s *StartRound) Hands() [common.NumPlayers][initHandSize]tile.Tile {
	return s.hands
}
