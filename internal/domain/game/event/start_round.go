package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

const (
	initHandSize = 13
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
