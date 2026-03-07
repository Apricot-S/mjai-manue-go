package round

import (
	"slices"

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
	if s.RoundNumber() == 4 {
		return s.RoundWind().Next(), 1
	}
	return s.RoundWind(), s.RoundNumber() + 1
}

func (s *State) Doras() tile.Tiles {
	doras := make([]tile.Tile, s.DoraIndicators().Len())
	for i := range doras {
		doras[i] = *s.doraIndicators[i].NextForDora()
	}
	return doras
}

func (s *State) Turn() float64 {
	return float64(NumInitWall-s.NumLeftTiles()) / float64(NumPlayers)
}

func (s *State) SeatWind(playerID id.ID) wind.Wind {
	return wind.Wind((playerID.Index()+1-s.RoundNumber()+4)%4 + 1)
}

func (s *State) VisibleTiles(playerID id.ID) tile.Tiles {
	var visibleTiles tile.Tiles

	for i := range NumPlayers {
		p := s.players[i]
		visibleTiles = slices.Concat(visibleTiles, p.River())

		for _, m := range p.Melds() {
			visibleTiles = slices.Concat(visibleTiles, m.ToTiles())
		}
	}

	var handTiles tile.Tiles
	if h, isVisible := s.Player(playerID).Hand(); isVisible {
		handTiles = h.ToTiles()
	}

	return slices.Concat(visibleTiles, s.DoraIndicators(), handTiles)
}

func (s *State) SafeTiles(playerID id.ID) tile.Tiles {
	p := s.Player(playerID)
	return slices.Concat(p.DiscardedTiles(), p.ExtraSafeTiles())
}
