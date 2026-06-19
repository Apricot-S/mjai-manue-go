package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func (s *State) LegalActions(playerSeat seat.Seat) ([]action.Action, error) {
	if actions, ok := s.legalActionsCache[playerSeat]; ok {
		return slices.Clone(actions), nil
	}
	if s.legalActionsCache == nil {
		s.legalActionsCache = map[seat.Seat][]action.Action{}
	}

	actions, err := s.calculateLegalActions(playerSeat)
	if err != nil {
		return nil, err
	}
	s.legalActionsCache[playerSeat] = actions
	return slices.Clone(actions), nil
}

func (s *State) calculateLegalActions(playerSeat seat.Seat) ([]action.Action, error) {
	if s.roundEnded {
		return nil, nil
	}
	if s.legalActionsSuppressed {
		return nil, nil
	}

	visiblePlayer, ok := s.players[playerSeat.Index()].(*player.VisiblePlayer)
	if !ok {
		return nil, fmt.Errorf("cannot list legal actions: player %d is invisible", playerSeat.Index())
	}

	if s.pendingDiscard == nil {
		// After open kan, the bot must choose a discard on the replacement draw,
		// even though the server publishes the dora event before the discard event.
		if s.pendingKanActor != nil && *s.pendingKanActor == playerSeat && s.kanProgress == noKanProgress && s.pendingDoraReveals > 0 {
			return s.legalActionsOnSelfDraw(playerSeat, visiblePlayer)
		}
		return s.legalActionsOnOtherDiscard(playerSeat, visiblePlayer)
	}
	if *s.pendingDiscard != playerSeat {
		return nil, nil
	}
	return s.legalActionsOnSelfDraw(playerSeat, visiblePlayer)
}
