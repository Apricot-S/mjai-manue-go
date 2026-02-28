package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/id"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
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
	NextRound() (wind.Wind, int)
	Doras() tile.Tiles
	SeatWind(playerID id.ID) wind.Wind
	VisibleTiles(playerID id.ID) tile.Tiles
	SafeTiles(playerID id.ID) tile.Tiles
}

type ActionStateViewer interface {
	CanAct(playerID id.ID) bool
}

type StateViewer interface {
	RawStateViewer
	DerivedStateViewer
	ActionStateViewer
}
