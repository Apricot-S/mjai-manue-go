package player

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type VisiblePlayer struct {
	hand                      hand.VisibleHand
	drawnTile                 *tile.Tile
	melds                     []meld.Meld
	river                     []tile.Tile
	discardedTiles            []tile.Tile
	extraSafeTiles            []tile.Tile
	riichiState               RiichiState
	riichiRiverIndex          int
	riichiDiscardedTilesIndex int
	isConcealed               bool
	swapCallTiles             []tile.Tile
	needsDeadWallDraw         bool
}

func NewVisiblePlayer(handTiles [initHandSize]tile.Tile) (*VisiblePlayer, error) {
	h, err := hand.NewVisibleHand(handTiles[:])
	if err != nil {
		return nil, err
	}

	return &VisiblePlayer{
		hand:                      *h,
		drawnTile:                 nil,
		melds:                     make([]meld.Meld, 0, maxNumMelds),
		river:                     make([]tile.Tile, 0, maxNumRiver),
		discardedTiles:            make([]tile.Tile, 0, maxNumDiscardedTiles),
		extraSafeTiles:            make([]tile.Tile, 0, 3),
		riichiState:               NotRiichi,
		riichiRiverIndex:          -1,
		riichiDiscardedTilesIndex: -1,
		isConcealed:               true,
		swapCallTiles:             nil,
		needsDeadWallDraw:         false,
	}, nil
}

func (p *VisiblePlayer) Hand() (*hand.VisibleHand, bool) {
	return &p.hand, true
}

func (p *VisiblePlayer) HandTiles() []tile.Tile {
	ts := tile.Tiles(p.hand.ToTiles())
	sort.Sort(ts)
	return ts
}

func (p *VisiblePlayer) DrawnTile() *tile.Tile {
	return p.drawnTile
}

func (p *VisiblePlayer) Melds() []meld.Meld {
	return p.melds
}

func (p *VisiblePlayer) River() []tile.Tile {
	return p.river
}

func (p *VisiblePlayer) DiscardedTiles() []tile.Tile {
	return p.discardedTiles
}

func (p *VisiblePlayer) ExtraSafeTiles() []tile.Tile {
	return p.extraSafeTiles
}

func (p *VisiblePlayer) RiichiState() RiichiState {
	return p.riichiState
}

func (p *VisiblePlayer) RiichiRiverIndex() int {
	return p.riichiRiverIndex
}

func (p *VisiblePlayer) RiichiDiscardedTilesIndex() int {
	return p.riichiDiscardedTilesIndex
}

func (p *VisiblePlayer) CanDiscard() bool {
	return !p.needsDeadWallDraw && p.drawnTile != nil || p.swapCallTiles != nil
}

func (p *VisiblePlayer) CanChiiPonKan() bool {
	return !p.needsDeadWallDraw && !p.CanDiscard() && len(p.Melds()) < 4
}

func (p *VisiblePlayer) IsConcealed() bool {
	return p.isConcealed
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

		for _, s := range p.swapCallTiles {
			if t.HasSameSymbol(&s) {
				return fmt.Errorf("cannot Discard: tile %s is forbidden due to swap-call", t)
			}
		}

		newHand, err := p.hand.Discard(&t)
		if err != nil {
			return err
		}

		if p.drawnTile != nil {
			if newHand, err = newHand.Draw(p.drawnTile); err != nil {
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
	allSwap := true
	for _, rt := range remaining.Distinct(nil) {
		found := false
		for _, s := range swapCallTiles {
			if rt.HasSameSymbol(&s) {
				found = true
				break
			}
		}
		if !found {
			allSwap = false
			break
		}
	}
	if allSwap {
		return fmt.Errorf("cannot Chii: remaining hand would contain only swap-call tiles")
	}

	p.hand = *h
	p.melds = append(p.melds, &chii)
	p.isConcealed = false
	p.swapCallTiles = swapCallTiles
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

	h, err := p.hand.Draw(p.drawnTile)
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
}

func (p *VisiblePlayer) TakeFromRiver(t tile.Tile) error {
	numRiver := len(p.river)

	if t != p.river[numRiver-1] {
		return fmt.Errorf("cannot take tile %s; last river tile is %s", t, p.river[numRiver-1])
	}

	p.river = slices.Delete(p.river, numRiver-1, numRiver)
	return nil
}
