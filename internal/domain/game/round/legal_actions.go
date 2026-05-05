package round

import (
	"fmt"
	"maps"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/action"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

const maxNumActions = 13 + 1 + 1 // discard + riichi + win
const maxKanCandidates = 3

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

	visiblePlayer, ok := s.players[playerSeat.Index()].(*player.VisiblePlayer)
	if !ok {
		return nil, fmt.Errorf("cannot list legal actions: player %d is invisible", playerSeat.Index())
	}

	if s.pendingDiscard == nil || *s.pendingDiscard != playerSeat {
		return nil, nil
	}

	return s.legalActionsOnSelfDraw(playerSeat, visiblePlayer)
}

func (s *State) legalActionsOnSelfDraw(playerSeat seat.Seat, p *player.VisiblePlayer) ([]action.Action, error) {
	if !p.CanDiscard() {
		return nil, fmt.Errorf("cannot list discard actions: player %d cannot discard", playerSeat.Index())
	}

	actions := make([]action.Action, 0, maxNumActions)
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
		concealedKans, err := s.legalConcealedKanActions(playerSeat, p)
		if err != nil {
			return nil, err
		}
		actions = append(actions, concealedKans...)
		return actions, nil
	}

	for _, handTile := range tile.Tiles(p.HandTiles()).Distinct(nil) {
		if isSwapCallTile(handTile, p.SwapCallTiles()) {
			continue
		}
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

	if s.canRiichi(p) {
		actions = append(actions, action.NewRiichi(playerSeat))
	}

	if s.canDeclareKyushukyuhai(playerSeat, p) {
		actions = append(actions, action.NewKyushukyuhai(playerSeat))
	}

	promotedKans, err := s.legalPromotedKanActions(playerSeat, p)
	if err != nil {
		return nil, err
	}
	actions = append(actions, promotedKans...)

	concealedKans, err := s.legalConcealedKanActions(playerSeat, p)
	if err != nil {
		return nil, err
	}
	actions = append(actions, concealedKans...)

	return actions, nil
}

func (s *State) canDeclareKyushukyuhai(playerSeat seat.Seat, p *player.VisiblePlayer) bool {
	if !s.canKyushukyuhai[playerSeat.Index()] {
		return false
	}

	handBeforeDeclare, ok := p.Hand()
	if !ok {
		return false
	}
	drawnTile := p.DrawnTile()
	if drawnTile == nil {
		return false
	}
	handAfterDraw, err := handBeforeDeclare.Draw(drawnTile)
	if err != nil {
		return false
	}

	tc34 := handAfterDraw.ToTileCounts34()
	numYaochuTypes := 0
	for _, id := range tile.YaochuhaiIDs {
		if tc34[id] > 0 {
			numYaochuTypes++
		}
	}
	return numYaochuTypes >= 9
}

func (s *State) legalConcealedKanActions(playerSeat seat.Seat, p *player.VisiblePlayer) ([]action.Action, error) {
	if s.numKans >= maxNumKan || s.numLeftTiles <= 0 {
		return nil, nil
	}
	if p.RiichiState() == player.RiichiDeclared {
		return nil, nil
	}

	handBeforeKan, ok := p.Hand()
	if !ok {
		return nil, nil
	}
	drawnTile := p.DrawnTile()
	if drawnTile == nil {
		return nil, nil
	}

	handAfterDraw, err := handBeforeKan.Draw(drawnTile)
	if err != nil {
		return nil, nil
	}

	actions := make([]action.Action, 0, maxKanCandidates)
	for id, count := range handAfterDraw.ToTileCounts34() {
		if count != 4 {
			continue
		}
		candidate := *tile.MustTileFromID(id)
		consumed := concealedKanConsumedTiles(candidate)
		if p.RiichiState() == player.RiichiAccepted && !canConcealedKanAfterRiichi(handBeforeKan, *drawnTile, consumed) {
			continue
		}

		k, err := meld.NewConcealedKan(consumed)
		if err != nil {
			continue
		}
		a, err := action.NewConcealedKan(playerSeat, [4]tile.Tile(k.Consumed()))
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func concealedKanConsumedTiles(candidate tile.Tile) [4]tile.Tile {
	return [4]tile.Tile{candidate, candidate, candidate, candidate.AddRed()}
}

func canConcealedKanAfterRiichi(handBeforeKan *hand.VisibleHand, drawnTile tile.Tile, consumed [4]tile.Tile) bool {
	if !drawnTile.HasSameSymbol(&consumed[0]) {
		return false
	}
	drawnTileID34 := drawnTile.RemoveRed().ID()
	tc34AfterKan := *handBeforeKan.ToTileCounts34()
	if tc34AfterKan[drawnTileID34] != 3 {
		return false
	}
	tc34AfterKan[drawnTileID34] = 0
	handAfterKan := hand.MustVisibleHand(tc34AfterKan.ToTiles())
	return maps.Equal(waitsFor(handBeforeKan), waitsFor(handAfterKan))
}

type waitSet map[int]struct{}

func waitsFor(h *hand.VisibleHand) waitSet {
	waits := waitSet{}
	for id := range tile.NumTileType34 {
		waitTile := tile.MustTileFromID(id)
		handWithWait, err := h.Draw(waitTile)
		if err != nil {
			continue
		}
		if service.IsWinningForm(handWithWait) {
			waits[id] = struct{}{}
		}
	}
	return waits
}

func (s *State) legalPromotedKanActions(playerSeat seat.Seat, p *player.VisiblePlayer) ([]action.Action, error) {
	if s.numKans >= maxNumKan || s.numLeftTiles <= 0 {
		return nil, nil
	}

	addedTiles := tile.Tiles(p.HandTiles())
	if drawnTile := p.DrawnTile(); drawnTile != nil {
		addedTiles = append(addedTiles, *drawnTile)
	}
	addedTiles = addedTiles.Distinct(nil)

	actions := make([]action.Action, 0, maxKanCandidates)
	for _, m := range p.Melds() {
		pon, ok := m.(*meld.Pon)
		if !ok {
			continue
		}

		for _, added := range addedTiles {
			k, err := meld.NewPromotedKan(*pon.Taken(), [2]tile.Tile(pon.Consumed()), added, *pon.Target())
			if err != nil {
				continue
			}
			a, err := action.NewPromotedKan(playerSeat, *k.Added(), [3]tile.Tile(pon.ToTiles()))
			if err != nil {
				return nil, err
			}
			actions = append(actions, a)
		}
	}
	return actions, nil
}

func (s *State) canRiichi(p *player.VisiblePlayer) bool {
	if p.RiichiState() != player.NotRiichi {
		return false
	}
	if p.DrawnTile() == nil {
		return false
	}
	if !p.IsConcealed() {
		return false
	}
	if s.numLeftTiles < common.NumPlayers {
		return false
	}

	handBeforeRiichi, ok := p.Hand()
	if !ok {
		return false
	}
	handAfterDraw, err := handBeforeRiichi.Draw(p.DrawnTile())
	if err != nil {
		return false
	}
	return service.IsTenpaiAll(handAfterDraw)
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

func isSwapCallTile(t tile.Tile, swapCallTiles []tile.Tile) bool {
	return slices.ContainsFunc(swapCallTiles, func(s tile.Tile) bool {
		return t.HasSameSymbol(&s)
	})
}
