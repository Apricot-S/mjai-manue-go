package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

const (
	MaxNumDoraIndicators = 5
	NumInitWall          = tile.NumTileType34*4 - 13*common.NumPlayers - 14
	FinalTurn            = float64(NumInitWall) / float64(common.NumPlayers)

	maxNumKan = 4
	// Indicates that no action has been taken by anyone.
	noActor                 = -1
	kanPlayerStatusMultiple = 4
)

type State struct {
	roundWind      wind.Wind
	roundNumber    int
	honba          int
	riichiDeposit  int
	scores         [common.NumPlayers]int
	dealer         seat.Seat
	startingDealer seat.Seat
	doraIndicators tile.Tiles
	numLeftTiles   int
	players        [common.NumPlayers]player.Player
}
