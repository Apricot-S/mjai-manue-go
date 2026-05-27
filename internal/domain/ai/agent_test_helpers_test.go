package ai

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

type stubPlayerViewer struct {
	hand                      *hand.VisibleHand
	riichiState               player.RiichiState
	drawnTile                 *tile.Tile
	discardedTiles            []tile.Tile
	riichiDiscardedTilesIndex int
	hasRiichiDiscardIndex     bool
}

func (p stubPlayerViewer) Hand() (*hand.VisibleHand, bool) {
	if p.hand == nil {
		return nil, false
	}
	return p.hand, true
}
func (p stubPlayerViewer) HandTiles() []tile.Tile          { return nil }
func (p stubPlayerViewer) DrawnTile() *tile.Tile           { return p.drawnTile }
func (p stubPlayerViewer) Melds() []meld.Meld              { return nil }
func (p stubPlayerViewer) River() []tile.Tile              { return p.discardedTiles }
func (p stubPlayerViewer) DiscardedTiles() []tile.Tile     { return p.discardedTiles }
func (p stubPlayerViewer) ExtraSafeTiles() []tile.Tile     { return nil }
func (p stubPlayerViewer) IsFuriten() bool                 { return false }
func (p stubPlayerViewer) CanRonBy(*tile.Tile) bool        { return true }
func (p stubPlayerViewer) RiichiState() player.RiichiState { return p.riichiState }
func (p stubPlayerViewer) RiichiRiverIndex() int           { return p.riichiIndex() }
func (p stubPlayerViewer) RiichiDiscardedTilesIndex() int  { return p.riichiIndex() }
func (p stubPlayerViewer) CanDiscard() bool                { return p.drawnTile != nil }
func (p stubPlayerViewer) CanChiiPonKan() bool             { return p.drawnTile == nil }
func (p stubPlayerViewer) IsConcealed() bool               { return true }
func (p stubPlayerViewer) SwapCallTiles() []tile.Tile      { return nil }

func (p stubPlayerViewer) riichiIndex() int {
	if !p.hasRiichiDiscardIndex {
		return -1
	}
	return p.riichiDiscardedTilesIndex
}

type stubWinEstimateStateViewer struct {
	turn         float64
	visibleTiles []tile.Tile
}

func (s stubWinEstimateStateViewer) VisibleTiles(seat.Seat) tile.Tiles {
	return s.visibleTiles
}

func (s stubWinEstimateStateViewer) Turn() float64 {
	return s.turn
}

type stubCandidateEvaluationStateViewer struct {
	turn            float64
	roundWind       wind.Wind
	roundNumber     int
	honba           int
	riichiDeposit   int
	seatWinds       [common.NumPlayers]wind.Wind
	dealer          seat.Seat
	scores          [common.NumPlayers]int
	startingSeat    seat.Seat
	visibleTiles    tile.Tiles
	safeTiles       [common.NumPlayers]tile.Tiles
	doras           tile.Tiles
	doraIndicators  tile.Tiles
	players         [common.NumPlayers]player.PlayerViewer
	nextRoundWind   wind.Wind
	nextRoundNumber int
	numLeftTiles    int
}

func stubStateWithSelf(self player.PlayerViewer) stubCandidateEvaluationStateViewer {
	selfSeat := seat.MustSeat(0)
	return stubCandidateEvaluationStateViewer{
		turn:         0,
		roundWind:    wind.East,
		roundNumber:  1,
		seatWinds:    [common.NumPlayers]wind.Wind{wind.East, wind.South, wind.West, wind.North},
		dealer:       selfSeat,
		scores:       [common.NumPlayers]int{25000, 25000, 25000, 25000},
		startingSeat: selfSeat,
		players: [common.NumPlayers]player.PlayerViewer{
			self,
			stubPlayerViewer{},
			stubPlayerViewer{},
			stubPlayerViewer{},
		},
		nextRoundWind:   wind.East,
		nextRoundNumber: 1,
		numLeftTiles:    round.NumInitWall,
	}
}

func (s stubCandidateEvaluationStateViewer) RoundWind() wind.Wind {
	return s.roundWind
}

func (s stubCandidateEvaluationStateViewer) RoundNumber() int {
	return s.roundNumber
}

func (s stubCandidateEvaluationStateViewer) Honba() int {
	return s.honba
}

func (s stubCandidateEvaluationStateViewer) RiichiDeposit() int {
	return s.riichiDeposit
}

func (s stubCandidateEvaluationStateViewer) SeatWind(playerSeat seat.Seat) wind.Wind {
	return s.seatWinds[playerSeat.Index()]
}

func (s stubCandidateEvaluationStateViewer) Doras() tile.Tiles {
	return s.doras
}

func (s stubCandidateEvaluationStateViewer) DoraIndicators() tile.Tiles {
	return s.doraIndicators
}

func (s stubCandidateEvaluationStateViewer) NumLeftTiles() int {
	return s.numLeftTiles
}

func (s stubCandidateEvaluationStateViewer) VisibleTiles(seat.Seat) tile.Tiles {
	return s.visibleTiles
}

func (s stubCandidateEvaluationStateViewer) SafeTiles(playerSeat seat.Seat) tile.Tiles {
	return s.safeTiles[playerSeat.Index()]
}

func (s stubCandidateEvaluationStateViewer) Player(playerSeat seat.Seat) player.PlayerViewer {
	return s.players[playerSeat.Index()]
}

func (s stubCandidateEvaluationStateViewer) Turn() float64 {
	return s.turn
}

func (s stubCandidateEvaluationStateViewer) Dealer() seat.Seat {
	return s.dealer
}

func (s stubCandidateEvaluationStateViewer) NextRound() (wind.Wind, int) {
	return s.nextRoundWind, s.nextRoundNumber
}

func (s stubCandidateEvaluationStateViewer) Scores() [common.NumPlayers]int {
	return s.scores
}

func (s stubCandidateEvaluationStateViewer) StartingDealer() seat.Seat {
	return s.startingSeat
}
