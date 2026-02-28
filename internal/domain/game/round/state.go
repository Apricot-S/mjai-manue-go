package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/id"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

const (
	NumPlayers           = 4
	MaxNumDoraIndicators = 5
	NumInitWall          = tile.NumTileType34*4 - 13*NumPlayers - 14
	FinalTurn            = float64(NumInitWall) / float64(NumPlayers)

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
	scores         [NumPlayers]int
	dealer         id.ID
	startingDealer id.ID
	doraIndicators tile.Tiles
	numLeftTiles   int
	players        [NumPlayers]player.Player
}
