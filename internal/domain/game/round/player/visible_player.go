package player

import (
	"fmt"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/meld"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type VisiblePlayer struct {
	hand      hand.VisibleHand
	drawnTile *tile.Tile
	melds     []meld.Meld
	river     []tile.Tile
}

func NewVisiblePlayer(handTiles []tile.Tile) (*VisiblePlayer, error) {
	if len(handTiles) != initHandSize {
		return nil, fmt.Errorf("invalid number of hand tiles: got %d, want %d", len(handTiles), initHandSize)
	}

	h, err := hand.NewVisibleHand(handTiles)
	if err != nil {
		return nil, err
	}

	return &VisiblePlayer{
		hand:      *h,
		drawnTile: nil,
		melds:     make([]meld.Meld, 0, maxNumMelds),
		river:     make([]tile.Tile, 0, maxNumRiver),
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
