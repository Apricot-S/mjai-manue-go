package player

import (
	"fmt"
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type InvisiblePlayer struct {
	hand hand.InvisibleHand
	commonPlayerState
}

var initInvisibleHandTiles tile.Tiles = tile.Tiles{
	tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?"),
	tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?"),
	tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?"), tile.MustTileFromCode("?"),
	tile.MustTileFromCode("?"),
}
var initInvisibleHand hand.InvisibleHand = *hand.MustInvisibleHand(initInvisibleHandTiles)

func NewInvisiblePlayer() *InvisiblePlayer {
	return &InvisiblePlayer{
		hand:              initInvisibleHand,
		commonPlayerState: newCommonPlayerState(),
	}
}

func (p *InvisiblePlayer) Hand() (*hand.VisibleHand, bool) {
	return nil, false
}

func (p *InvisiblePlayer) HandTiles() []tile.Tile {
	return p.hand.ToTiles()
}

func (p *InvisiblePlayer) IsFuriten() bool {
	return false
}

func (p *InvisiblePlayer) CanRonBy(winningTile *tile.Tile) bool {
	if winningTile == nil {
		// When the event omits the winning tile, state transition validation
		// cannot verify the exact tile, so preserve the existing permissive behavior.
		return true
	}

	if tile.Tiles(p.discardedTiles).ContainsSameSymbol(*winningTile) {
		return false
	}
	if tile.Tiles(p.extraSafeTiles).ContainsSameSymbol(*winningTile) {
		return false
	}
	return true
}

func (p *InvisiblePlayer) Draw(t tile.Tile) error {
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

func (p *InvisiblePlayer) Discard(t tile.Tile, tsumogiri bool) error {
	if !p.CanDiscard() {
		return fmt.Errorf("cannot Discard: player is not in a discardable state")
	}
	if t.IsUnknown() {
		return fmt.Errorf("invisible player cannot discard an unknown tile")
	}

	if !tsumogiri {
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

		p.hand = *newHand
	}

	if p.riichiState != RiichiAccepted {
		p.extraSafeTiles = make([]tile.Tile, 0, 3)
	}

	p.drawnTile = nil
	p.river = append(p.river, t)
	p.discardedTiles = append(p.discardedTiles, t)
	p.swapCallTiles = nil
	return nil
}

func (p *InvisiblePlayer) Chii(chii meld.Chii) error {
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

	p.hand = *h
	p.melds = append(p.melds, &chii)
	p.isConcealed = false
	p.swapCallTiles = chii.SwapCallTiles()
	return nil
}

func (p *InvisiblePlayer) Pon(pon meld.Pon) error {
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
	return nil
}

func (p *InvisiblePlayer) CalledKan(kan meld.CalledKan) error {
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
	return nil
}

func (p *InvisiblePlayer) ConcealedKan(kan meld.ConcealedKan) error {
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
	return nil
}

func (p *InvisiblePlayer) PromotedKan(kan meld.PromotedKan) error {
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
	return nil
}

func (p *InvisiblePlayer) Riichi() error {
	if p.riichiState != NotRiichi {
		return fmt.Errorf("cannot Riichi: player is already in riichi state (%v)", p.riichiState)
	}
	if p.drawnTile == nil {
		return fmt.Errorf("cannot Riichi: player is not in a discardable state")
	}
	if !p.isConcealed {
		return fmt.Errorf("cannot Riichi: player hand is not concealed")
	}

	p.riichiState = RiichiDeclared
	return nil
}

func (p *InvisiblePlayer) RiichiAccepted() error {
	if p.riichiState != RiichiDeclared {
		return fmt.Errorf("Riichi cannot be accepted: invalid state (%v)", p.riichiState)
	}

	p.riichiState = RiichiAccepted
	p.riichiRiverIndex = len(p.River()) - 1
	p.riichiDiscardedTilesIndex = len(p.DiscardedTiles()) - 1
	return nil
}

func (p *InvisiblePlayer) AddExtraSafeTiles(t tile.Tile) {
	if t.IsUnknown() {
		panic("cannot add an unknown tile to extraSafeTiles")
	}

	p.extraSafeTiles = append(p.extraSafeTiles, t)
}
