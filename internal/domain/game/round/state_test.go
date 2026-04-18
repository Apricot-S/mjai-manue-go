package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func NewStateForTest(
	roundWind wind.Wind,
	roundNumber int,
	honba int,
	riichiDeposit int,
	scores [common.NumPlayers]int,
	dealer seat.Seat,
	startingDealer seat.Seat,
	doraIndicators tile.Tiles,
	numLeftTiles int,
	players [common.NumPlayers]player.Player,
) State {
	return State{
		roundWind,
		roundNumber,
		honba,
		riichiDeposit,
		scores,
		dealer,
		startingDealer,
		doraIndicators,
		numLeftTiles,
		players,
	}
}
