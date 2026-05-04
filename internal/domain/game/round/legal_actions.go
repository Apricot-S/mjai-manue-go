package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
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
	if s.pendingDiscard == nil || *s.pendingDiscard != playerSeat {
		return nil, nil
	}

	return s.legalDiscardActions(playerSeat)
}

func (s *State) legalDiscardActions(playerSeat seat.Seat) ([]action.Action, error) {
	p := s.players[playerSeat.Index()]
	if !p.CanDiscard() {
		return nil, fmt.Errorf("cannot list discard actions: player %d cannot discard", playerSeat.Index())
	}

	actions := make([]action.Action, 0, 14)
	addDiscard := func(discardedTile tile.Tile, tsumogiri bool) error {
		a, err := action.NewDiscard(playerSeat, discardedTile, tsumogiri)
		if err != nil {
			return err
		}
		actions = append(actions, a)
		return nil
	}

	for _, handTile := range tile.Tiles(p.HandTiles()).Distinct(nil) {
		if err := addDiscard(handTile, false); err != nil {
			return nil, err
		}
	}

	if drawnTile := p.DrawnTile(); drawnTile != nil {
		if err := addDiscard(*drawnTile, true); err != nil {
			return nil, err
		}
	}

	return actions, nil
}
