package player

import (
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type VisiblePlayer struct {
	hand hand.VisibleHand
}

func NewVisiblePlayer(handTiles []tile.Tile) (*VisiblePlayer, error) {
	h, _ := hand.NewVisibleHand(handTiles)
	return &VisiblePlayer{hand: *h}, nil
}

func (p *VisiblePlayer) Hand() (*hand.VisibleHand, bool) {
	return &p.hand, true
}

func (p *VisiblePlayer) HandTiles() []tile.Tile {
	ts := tile.Tiles(p.hand.ToTiles())
	sort.Sort(tile.Tiles(p.hand.ToTiles()))
	return ts
}
