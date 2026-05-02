package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/wind"
)

const (
	MaxNumDoraIndicators = 5
	NumInitWall          = tile.NumTileType34*4 - 13*common.NumPlayers - 14
	FinalTurn            = float64(NumInitWall) / float64(common.NumPlayers)

	minRoundNumber = 1
	maxRoundNumber = 4
	maxNumKan      = 4
)

type kanProgress int

const (
	noKanProgress kanProgress = iota
	// waitingReplacementBeforeDora is used after daiminkan/kakan and waits for the kan actor's replacement tile draw.
	waitingReplacementBeforeDora
	// waitingDoraAfterReplacement is used after daiminkan/kakan replacement tile draw and waits for the dora reveal.
	waitingDoraAfterReplacement
	// waitingDoraBeforeReplacement is used after ankan and waits for the dora reveal before replacement tile draw.
	waitingDoraBeforeReplacement
	// waitingReplacementAfterDora is used after ankan dora reveal and waits for the kan actor's replacement tile draw.
	waitingReplacementAfterDora
)

type State struct {
	roundWind               wind.Wind
	roundNumber             int
	honba                   int
	riichiDeposit           int
	scores                  [common.NumPlayers]int
	dealer                  seat.Seat
	startingDealer          seat.Seat
	doraIndicators          tile.Tiles
	numLeftTiles            int
	numKans                 int
	kanProgress             kanProgress
	pendingKanActor         *seat.Seat
	nextDraw                seat.Seat
	pendingDiscard          *seat.Seat
	pendingRiichiAcceptance *seat.Seat
	lastActor               *seat.Seat
	players                 [common.NumPlayers]player.Player
}

func NewState(ev *event.StartRound, previousScores [common.NumPlayers]int) (*State, error) {
	if ev.RoundWind() < wind.East || wind.North < ev.RoundWind() {
		return nil, fmt.Errorf("invalid round wind: %v", ev.RoundWind())
	}
	if ev.RoundNumber() < minRoundNumber || maxRoundNumber < ev.RoundNumber() {
		return nil, fmt.Errorf("invalid round number: %d", ev.RoundNumber())
	}
	if ev.Honba() < 0 {
		return nil, fmt.Errorf("invalid honba: %d", ev.Honba())
	}
	if ev.RiichiDeposit() < 0 {
		return nil, fmt.Errorf("invalid riichi deposit: %d", ev.RiichiDeposit())
	}
	if ev.DoraIndicator().IsUnknown() {
		return nil, fmt.Errorf("invalid dora indicator: %v", ev.DoraIndicator())
	}

	s := &State{}

	s.roundWind = ev.RoundWind()
	s.roundNumber = ev.RoundNumber()
	s.honba = ev.Honba()
	s.riichiDeposit = ev.RiichiDeposit()
	s.dealer = ev.Dealer()
	s.startingDealer = *seat.MustSeat(0)
	s.doraIndicators = make(tile.Tiles, 0, MaxNumDoraIndicators)
	s.doraIndicators = append(s.doraIndicators, ev.DoraIndicator())
	s.numLeftTiles = NumInitWall
	s.nextDraw = ev.Dealer()

	if ev.Scores() != nil {
		s.scores = *ev.Scores()
	} else {
		s.scores = previousScores
	}

	for i, handTiles := range ev.Hands() {
		p, err := s.newPlayerFromHand(&handTiles)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize player %d: %w", i, err)
		}
		s.players[i] = p
	}

	return s, nil
}

func (s *State) newPlayerFromHand(handTiles *[common.InitHandSize]tile.Tile) (player.Player, error) {
	if isUnknownHand(handTiles) {
		return player.NewInvisiblePlayer(), nil
	}

	visiblePlayer, err := player.NewVisiblePlayer(*handTiles)
	if err != nil {
		return nil, err
	}
	return visiblePlayer, nil
}

func isUnknownHand(handTiles *[common.InitHandSize]tile.Tile) bool {
	return slices.IndexFunc(handTiles[:], func(t tile.Tile) bool {
		return !t.IsUnknown()
	}) == -1
}
