package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/id"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

func NewStateForTest(
	roundWind wind.Wind,
	roundNumber int,
	honba int,
	riichiDeposit int,
	scores [NumPlayers]int,
	dealer id.ID,
	startingDealer id.ID,
	doraIndicators tile.Tiles,
	numLeftTiles int,
	players [NumPlayers]player.Player,
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
