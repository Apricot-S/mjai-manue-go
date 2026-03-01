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
	Player(playerID id.ID) player.PlayerViewer
}

type DerivedStateViewer interface {
	NextRound() (wind.Wind, int)
	Doras() tile.Tiles
	Turn() float64
	SeatWind(playerID id.ID) wind.Wind
	VisibleTiles(playerID id.ID) tile.Tiles
	SafeTiles(playerID id.ID) tile.Tiles
}

type StateViewer interface {
	RawStateViewer
	DerivedStateViewer
}

func (s *State) RoundWind() wind.Wind {
	return s.roundWind
}

func (s *State) RoundNumber() int {
	return s.roundNumber
}

func (s *State) Honba() int {
	return s.honba
}

func (s *State) RiichiDeposit() int {
	return s.riichiDeposit
}

func (s *State) Scores() [NumPlayers]int {
	return s.scores
}

func (s *State) Dealer() id.ID {
	return s.dealer
}

func (s *State) StartingDealer() id.ID {
	return s.startingDealer
}

func (s *State) DoraIndicators() tile.Tiles {
	return s.doraIndicators
}

func (s *State) NumLeftTiles() int {
	return s.numLeftTiles
}

func (s *State) Player(playerID id.ID) player.PlayerViewer {
	return s.players[playerID.Index()]
}

func (s *State) NextRound() (wind.Wind, int) {
	return wind.East, 2
}
