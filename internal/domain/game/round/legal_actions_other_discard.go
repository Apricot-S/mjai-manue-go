package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

const maxNumActionsOnOtherDiscard = 1 + 1 + 5 // pon + daiminkan + up to 5 chii patterns with red fives

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

	pons, err := s.legalPonActions(playerSeat, p, *targetSeat, discardedTile)
	if err != nil {
		return nil, err
	}
	actions = append(actions, pons...)

	calledKans, err := s.legalCalledKanActions(playerSeat, p, *targetSeat, discardedTile)
	if err != nil {
		return nil, err
	}
	actions = append(actions, calledKans...)

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

func (s *State) legalPonActions(playerSeat seat.Seat, p *player.VisiblePlayer, targetSeat seat.Seat, taken tile.Tile) ([]action.Action, error) {
	if !p.CanChiiPonKan() || s.numLeftTiles <= 0 {
		return nil, nil
	}

	handBeforeCall, ok := p.Hand()
	if !ok {
		return nil, nil
	}
	consumedCandidates := ponConsumedCandidates(handBeforeCall.Count, taken)
	if len(consumedCandidates) == 0 {
		return nil, nil
	}

	actions := make([]action.Action, 0, len(consumedCandidates))
	for _, consumed := range consumedCandidates {
		a, err := action.NewPon(playerSeat, targetSeat, taken, consumed)
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (s *State) legalCalledKanActions(playerSeat seat.Seat, p *player.VisiblePlayer, targetSeat seat.Seat, taken tile.Tile) ([]action.Action, error) {
	if !p.CanChiiPonKan() || s.numKans >= maxNumKan || s.numLeftTiles <= 0 {
		return nil, nil
	}

	handBeforeCall, ok := p.Hand()
	if !ok {
		return nil, nil
	}
	consumed, ok := calledKanConsumedCandidate(handBeforeCall.Count, taken)
	if !ok {
		return nil, nil
	}

	a, err := action.NewCalledKan(playerSeat, targetSeat, taken, consumed)
	if err != nil {
		return nil, err
	}
	return []action.Action{a}, nil
}

func calledKanConsumedCandidate(count func(tile.Tile) int, taken tile.Tile) ([3]tile.Tile, bool) {
	normal := taken.RemoveRed()
	if !normal.IsSuits() || normal.Number() != 5 {
		if count(normal) < 3 {
			return [3]tile.Tile{}, false
		}
		return [3]tile.Tile{normal, normal, normal}, true
	}

	red := normal.AddRed()
	if taken.IsRed() {
		if count(normal) < 3 {
			return [3]tile.Tile{}, false
		}
		return [3]tile.Tile{normal, normal, normal}, true
	}
	if count(normal) < 2 || count(red) < 1 {
		return [3]tile.Tile{}, false
	}
	return [3]tile.Tile{normal, normal, red}, true
}

func ponConsumedCandidates(count func(tile.Tile) int, taken tile.Tile) [][2]tile.Tile {
	normal := taken.RemoveRed()
	if !normal.IsSuits() || normal.Number() != 5 {
		if count(normal) < 2 {
			return nil
		}
		return [][2]tile.Tile{{normal, normal}}
	}

	red := normal.AddRed()
	normalCount := count(normal)
	redCount := count(red)
	consumed := make([][2]tile.Tile, 0, 2)
	if normalCount >= 2 {
		consumed = append(consumed, [2]tile.Tile{normal, normal})
	}
	if !taken.IsRed() && normalCount >= 1 && redCount >= 1 {
		consumed = append(consumed, [2]tile.Tile{normal, red})
	}
	return consumed
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
