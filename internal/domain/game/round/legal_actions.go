package round

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
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

	if p.RiichiState() == player.RiichiAccepted {
		drawnTile := p.DrawnTile()
		if drawnTile == nil {
			return nil, fmt.Errorf("cannot list discard actions: riichi player %d has no drawn tile", playerSeat.Index())
		}
		if err := addDiscard(*drawnTile, true); err != nil {
			return nil, err
		}
		return actions, nil
	}

	for _, handTile := range tile.Tiles(p.HandTiles()).Distinct(nil) {
		if p.RiichiState() == player.RiichiDeclared && !canDiscardAsRiichiDeclarationTile(p, handTile, false) {
			continue
		}
		if err := addDiscard(handTile, false); err != nil {
			return nil, err
		}
	}

	if drawnTile := p.DrawnTile(); drawnTile != nil {
		if p.RiichiState() == player.RiichiDeclared && !canDiscardAsRiichiDeclarationTile(p, *drawnTile, true) {
			return actions, nil
		}
		if err := addDiscard(*drawnTile, true); err != nil {
			return nil, err
		}
	}

	return actions, nil
}

func canDiscardAsRiichiDeclarationTile(p player.Player, discardTile tile.Tile, tsumogiri bool) bool {
	handBeforeDiscard, ok := p.Hand()
	if !ok {
		return false
	}
	if tsumogiri {
		return service.IsTenpaiAll(handBeforeDiscard)
	}

	handAfterDiscard, err := handBeforeDiscard.Discard(&discardTile)
	if err != nil {
		return false
	}
	if drawnTile := p.DrawnTile(); drawnTile != nil {
		handAfterDiscard, err = handAfterDiscard.Draw(drawnTile)
		if err != nil {
			return false
		}
	}
	return service.IsTenpaiAll(handAfterDiscard)
}
