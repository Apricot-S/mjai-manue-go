package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
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
}
