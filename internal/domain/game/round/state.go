package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
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
	// Dora reveals are tracked separately because logs may reveal those kan dora after later consecutive kan events.
	waitingReplacementBeforeDora
	// waitingReplacementAfterDora is used after ankan and waits for pending dora reveals before the replacement tile draw.
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
	pendingDoraReveals      int
	pendingRobbedKanTile    *tile.Tile
	nextDraw                seat.Seat
	pendingDiscard          *seat.Seat
	pendingRiichiAcceptance *seat.Seat
	lastDrawWasReplacement  bool
	canKyushukyuhai         [common.NumPlayers]bool
	roundEnded              bool
	roundEndedByWin         bool
	winTarget               *seat.Seat
	winActors               [common.NumPlayers]bool
	lastActor               *seat.Seat
	players                 [common.NumPlayers]player.Player
	legalActionsCache       map[seat.Seat][]action.Action
}

func NewState(ev *event.StartRound, previousScores [common.NumPlayers]int) (*State, error) {
	roundWind := ev.RoundWind()
	roundNumber := ev.RoundNumber()
	honba := ev.Honba()
	riichiDeposit := ev.RiichiDeposit()
	dealer := ev.Dealer()
	doraIndicator := ev.DoraIndicator()

	if roundWind < wind.East || wind.North < roundWind {
		return nil, fmt.Errorf("invalid round wind: %v", roundWind)
	}
	if roundNumber < minRoundNumber || maxRoundNumber < roundNumber {
		return nil, fmt.Errorf("invalid round number: %d", roundNumber)
	}
	if honba < 0 {
		return nil, fmt.Errorf("invalid honba: %d", honba)
	}
	if riichiDeposit < 0 {
		return nil, fmt.Errorf("invalid riichi deposit: %d", riichiDeposit)
	}
	if doraIndicator.IsUnknown() {
		return nil, fmt.Errorf("invalid dora indicator: %v", doraIndicator)
	}

	s := &State{}

	s.roundWind = roundWind
	s.roundNumber = roundNumber
	s.honba = honba
	s.riichiDeposit = riichiDeposit
	s.dealer = dealer
	s.startingDealer = seat.MustSeat(0)
	s.doraIndicators = make(tile.Tiles, 0, MaxNumDoraIndicators)
	s.doraIndicators = append(s.doraIndicators, doraIndicator)
	s.numLeftTiles = NumInitWall
	s.nextDraw = dealer
	for i := range s.canKyushukyuhai {
		s.canKyushukyuhai[i] = true
	}

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
