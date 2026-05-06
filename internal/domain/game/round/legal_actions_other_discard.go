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
	return actions, nil
}

func (s *State) canWinByRon(playerSeat seat.Seat, p *player.VisiblePlayer, winningTile tile.Tile) bool {
	handBeforeWin, ok := p.Hand()
	if !ok {
		return false
	}
	return service.Has1Han(
		handBeforeWin,
		p.Melds(),
		&winningTile,
		s.roundWind,
		s.SeatWind(playerSeat),
		false,
		p.RiichiState() != player.NotRiichi,
		service.NoEvent,
	)
}
