package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/event"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type EventApplier interface {
	Apply(ev event.Event) error
}

func (s *State) Apply(ev event.Event) error {
	switch ev := ev.(type) {
	case *event.StartRound:
		return s.applyStartRound(ev)
	case *event.EndRound:
		return fmt.Errorf("unimplemented event: %T", ev)
	}
	return fmt.Errorf("unknown event: %T", ev)
}

func (s *State) applyStartRound(ev *event.StartRound) error {
	s.roundWind = ev.RoundWind()
	s.roundNumber = ev.RoundNumber()
	s.honba = ev.Honba()
	s.riichiDeposit = ev.RiichiDeposit()
	s.dealer = ev.Dealer()
	s.startingDealer = ev.StartingDealer()
	s.doraIndicators = tile.Tiles{ev.DoraIndicator()}
	s.numLeftTiles = NumInitWall

	if ev.Scores() != nil {
		s.scores = *ev.Scores()
	}

	for i, handTiles := range ev.Hands() {
		p, err := s.newPlayerFromHand(&handTiles)
		if err != nil {
			return fmt.Errorf("failed to initialize player %d: %w", i, err)
		}
		s.players[i] = p
	}

	return nil
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
