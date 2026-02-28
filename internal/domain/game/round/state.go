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

type RawStateViewer interface {
	RoundWind() wind.Wind
	RoundNumber() int
	Honba() int
	RiichiDeposit() int
	Scores() [NumPlayers]int
	Dealer() id.ID
	StartingDealer() id.ID
	DoraIndicators() tile.Tiles
	NumLeftTiles() int
	Turn() float64
	Players() *[NumPlayers]player.PlayerViewer
}

type DerivedStateViewer interface {
	Doras() tile.Tiles
	SeatWind(playerID id.ID) wind.Wind
	VisibleTiles(playerID id.ID) tile.Tiles
	SafeTiles(playerID id.ID) tile.Tiles
	NextRound() (wind.Wind, int)
}

type ActionStateViewer interface {
	CanAct(playerID id.ID) bool
}

type StateViewer interface {
	RawStateViewer
	DerivedStateViewer
	ActionStateViewer
}
