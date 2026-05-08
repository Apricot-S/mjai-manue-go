package player

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/common"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type VisiblePlayer struct {
	hand hand.VisibleHand
	commonPlayerState
	waits     service.WaitSet
	isFuriten bool
}

func NewVisiblePlayer(handTiles [common.InitHandSize]tile.Tile) (*VisiblePlayer, error) {
	h, err := hand.NewVisibleHand(handTiles[:])
	if err != nil {
		return nil, err
	}

	p := &VisiblePlayer{
		hand:              *h,
		commonPlayerState: newCommonPlayerState(),
	}
	p.updateWaits()
	return p, nil
}

func (p *VisiblePlayer) Hand() (*hand.VisibleHand, bool) {
	return &p.hand, true
}

func (p *VisiblePlayer) HandTiles() []tile.Tile {
	ts := tile.Tiles(p.hand.ToTiles())
	ts.Sort()
	return ts
}

func (p *VisiblePlayer) IsFuriten() bool {
	return p.isFuriten
}

func (p *VisiblePlayer) CanRonBy(winningTile *tile.Tile) bool {
	if winningTile == nil {
		// When the event omits the winning tile, state transition validation
		// cannot verify the exact tile, so preserve the existing permissive behavior.
		return true
	}
	return !p.isFuriten && p.waits.Has(*winningTile)
}

func (p *VisiblePlayer) Draw(t tile.Tile) error {
	if t.IsUnknown() {
		return fmt.Errorf("visible player cannot draw an unknown tile")
	}
	if p.CanDiscard() {
		return fmt.Errorf("cannot Draw: player is already in a discardable state")
	}
	if p.riichiState == RiichiDeclared {
		return fmt.Errorf("cannot Draw: while declaring Riichi")
	}

	p.drawnTile = &t
	p.needsDeadWallDraw = false
	return nil
}

func (p *VisiblePlayer) Discard(t tile.Tile, tsumogiri bool) error {
	if !p.CanDiscard() {
		return fmt.Errorf("cannot Discard: player is not in a discardable state")
	}

	if tsumogiri {
		if t != *p.drawnTile {
			return fmt.Errorf("cannot Discard: tsumogiri tile (%s) must equal the drawn tile (%s)", t, p.drawnTile)
		}
		if p.riichiState == RiichiDeclared && !service.IsTenpaiAll(&p.hand) {
			return fmt.Errorf("cannot Discard: player is in riichi and discarding %s would break tenpai", t)
		}
	} else {
		if p.riichiState == RiichiAccepted {
			return fmt.Errorf("cannot Discard: player has accepted riichi and cannot discard a tile from hand: %s", t)
		}

		if isSwapCallTile(t, p.swapCallTiles) {
			return fmt.Errorf("cannot Discard: tile %s is forbidden due to swap-call", t)
		}

		newHand, err := p.hand.Discard(t)
		if err != nil {
			return err
		}

		if p.drawnTile != nil {
			if newHand, err = newHand.Draw(*p.drawnTile); err != nil {
				return err
			}
		}

		if p.riichiState == RiichiDeclared && !service.IsTenpaiAll(newHand) {
			return fmt.Errorf("cannot Discard: player is in riichi and discarding %s would break tenpai", t)
		}

		p.hand = *newHand
	}

	if p.riichiState != RiichiAccepted {
		p.extraSafeTiles = make([]tile.Tile, 0, 3)
	}

	p.drawnTile = nil
	p.river = append(p.river, t)
	p.discardedTiles = append(p.discardedTiles, t)
	p.swapCallTiles = nil
	p.updateWaits()
	p.updateFuritenAfterDiscard()
	return nil
}

func (p *VisiblePlayer) Chii(chii meld.Chii) error {
	if p.riichiState != NotRiichi {
		return fmt.Errorf("cannot Chii: player is already in riichi state (%v)", p.riichiState)
	}
	if !p.CanChiiPonKan() {
		return fmt.Errorf("cannot Chii: player is in a discardable state")
	}

	h, err := p.hand.Call(&chii)
	if err != nil {
		return err
	}

	// If the only tiles remaining after chii are swap-call tiles, chii is not allowed.
	swapCallTiles := chii.SwapCallTiles()
	remaining := tile.Tiles(h.ToTiles())
	hasNonSwapCallTile := slices.ContainsFunc(remaining.Distinct(nil), func(rt tile.Tile) bool {
		return !isSwapCallTile(rt, swapCallTiles)
	})
	if !hasNonSwapCallTile {
		return fmt.Errorf("cannot Chii: remaining hand would contain only swap-call tiles")
	}

	p.hand = *h
	p.melds = append(p.melds, &chii)
	p.isConcealed = false
	p.swapCallTiles = swapCallTiles
	p.updateWaits()
	return nil
}

func (p *VisiblePlayer) Pon(pon meld.Pon) error {
	if p.riichiState != NotRiichi {
		return fmt.Errorf("cannot Pon: player is already in riichi state (%v)", p.riichiState)
	}
	if !p.CanChiiPonKan() {
		return fmt.Errorf("cannot Pon: player is in a discardable state")
	}

	h, err := p.hand.Call(&pon)
	if err != nil {
		return err
	}

	p.hand = *h
	p.melds = append(p.melds, &pon)
	p.isConcealed = false
	p.swapCallTiles = pon.SwapCallTiles()
	p.updateWaits()
	return nil
}

func (p *VisiblePlayer) CalledKan(kan meld.CalledKan) error {
	if p.riichiState != NotRiichi {
		return fmt.Errorf("cannot CalledKan: player is already in riichi state (%v)", p.riichiState)
	}
	if !p.CanChiiPonKan() {
		return fmt.Errorf("cannot CalledKan: player is in a discardable state")
	}

	h, err := p.hand.Call(&kan)
	if err != nil {
		return err
	}

	p.hand = *h
	p.melds = append(p.melds, &kan)
	p.isConcealed = false
	p.needsDeadWallDraw = true
	p.updateWaits()
	return nil
}

func (p *VisiblePlayer) ConcealedKan(kan meld.ConcealedKan) error {
	if !p.CanDiscard() {
		return fmt.Errorf("cannot ConcealedKan: player is not in a discardable state")
	}

	newHand, err := p.hand.Draw(*p.drawnTile)
	if err != nil {
		return err
	}

	h, err := newHand.Call(&kan)
	if err != nil {
		return err
	}

	p.hand = *h
	p.drawnTile = nil
	p.melds = append(p.melds, &kan)
	p.needsDeadWallDraw = true
	p.updateWaits()
	return nil
}

func (p *VisiblePlayer) PromotedKan(kan meld.PromotedKan) error {
	if !p.CanDiscard() {
		return fmt.Errorf("cannot PromotedKan: player is not in a discardable state")
	}

	ponIndex := slices.IndexFunc(p.melds, func(m meld.Meld) bool {
		pon, isPon := m.(*meld.Pon)
		if !isPon {
			return false
		}
		return pon.Taken() == kan.Taken()
	})
	if ponIndex == -1 {
		return fmt.Errorf("cannot PromotedKan: failed to find pon for promoted kan: %v", p.melds)
	}

	newHand, err := p.hand.Draw(*p.drawnTile)
	if err != nil {
		return err
	}

	h, err := newHand.Call(&kan)
	if err != nil {
		return err
	}

	p.hand = *h
	p.drawnTile = nil
	p.melds[ponIndex] = &kan
	p.needsDeadWallDraw = true
	p.updateWaits()
	return nil
}

func (p *VisiblePlayer) Riichi() error {
	if p.riichiState != NotRiichi {
		return fmt.Errorf("cannot Riichi: player is already in riichi state (%v)", p.riichiState)
	}
	if p.drawnTile == nil {
		return fmt.Errorf("cannot Riichi: player is not in a discardable state")
	}
	if !p.isConcealed {
		return fmt.Errorf("cannot Riichi: player hand is not concealed")
	}

	h, err := p.hand.Draw(*p.drawnTile)
	if err != nil {
		return err
	}
	if !service.IsTenpaiAll(h) {
		return fmt.Errorf("cannot Riichi: player is not tenpai")
	}

	p.riichiState = RiichiDeclared
	return nil
}

func (p *VisiblePlayer) RiichiAccepted() error {
	if p.riichiState != RiichiDeclared {
		return fmt.Errorf("Riichi cannot be accepted: invalid state (%v)", p.riichiState)
	}

	p.riichiState = RiichiAccepted
	p.riichiRiverIndex = len(p.River()) - 1
	p.riichiDiscardedTilesIndex = len(p.DiscardedTiles()) - 1
	return nil
}

func (p *VisiblePlayer) AddExtraSafeTiles(t tile.Tile) {
	if t.IsUnknown() {
		panic("cannot add an unknown tile to extraSafeTiles")
	}

	p.extraSafeTiles = append(p.extraSafeTiles, t)
	if p.waits.Has(t) {
		p.isFuriten = true
	}
}

func (p *VisiblePlayer) updateWaits() {
	p.waits = service.WaitsFor(&p.hand)
}

func (p *VisiblePlayer) updateFuritenAfterDiscard() {
	riverFuriten := slices.ContainsFunc(p.discardedTiles, p.waits.Has)
	if p.riichiState == RiichiAccepted {
		p.isFuriten = p.isFuriten || riverFuriten
		return
	}
	p.isFuriten = riverFuriten
}
