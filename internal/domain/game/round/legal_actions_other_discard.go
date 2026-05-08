package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

const maxNumActionsOnOtherDiscard = 1 + 5 // ron + up to 5 chii patterns with red fives

func (s *State) legalActionsOnOtherDiscard(playerSeat seat.Seat, p *player.VisiblePlayer) ([]action.Action, error) {
	if s.pendingRobbedKanTile != nil && s.pendingKanActor != nil && *s.pendingKanActor != playerSeat {
		return s.legalActionsOnRobbingKan(playerSeat, p, *s.pendingKanActor, *s.pendingRobbedKanTile)
	}

	targetSeat := s.lastActor
	if targetSeat == nil || *targetSeat == playerSeat {
		return nil, nil
	}
	target := s.players[targetSeat.Index()]
	river := target.River()
	if len(river) == 0 {
		return nil, nil
	}
	discardedTile := river[len(river)-1]

	actions := make([]action.Action, 0, maxNumActionsOnOtherDiscard)
	if s.canWinByRon(playerSeat, p, discardedTile) {
		a, err := action.NewWin(playerSeat, *targetSeat, discardedTile)
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}

	if len(actions) > 0 {
		actions = append(actions, action.NewPass(playerSeat))
	}
	return actions, nil
}

func (s *State) legalActionsOnRobbingKan(playerSeat seat.Seat, p *player.VisiblePlayer, targetSeat seat.Seat, winningTile tile.Tile) ([]action.Action, error) {
	actions := make([]action.Action, 0, 2)
	if s.canWinByRon(playerSeat, p, winningTile) {
		a, err := action.NewWin(playerSeat, targetSeat, winningTile)
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	if len(actions) > 0 {
		actions = append(actions, action.NewPass(playerSeat))
	}
	return actions, nil
}

func (s *State) canWinByRon(playerSeat seat.Seat, p *player.VisiblePlayer, winningTile tile.Tile) bool {
	if p.IsFuriten() {
		return false
	}

	handBeforeWin, ok := p.Hand()
	if !ok {
		return false
	}
	return service.Has1Han(
		handBeforeWin,
		p.Melds(),
		winningTile,
		s.roundWind,
		s.SeatWind(playerSeat),
		false,
		p.RiichiState() != player.NotRiichi,
		s.ronWinEvent(),
	)
}

func (s *State) ronWinEvent() service.WinEvent {
	if s.pendingRobbedKanTile != nil {
		return service.RobbingAKan
	}
	if s.numLeftTiles == 0 {
		return service.LastTile
	}
	return service.NoEvent
}
