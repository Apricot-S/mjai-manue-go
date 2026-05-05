package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
)

func (s *State) legalActionsOnOtherDiscard(playerSeat seat.Seat, p *player.VisiblePlayer) ([]action.Action, error) {
	return nil, nil
}
