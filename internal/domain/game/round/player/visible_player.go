package player

import (
	"fmt"
	"sort"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type VisiblePlayer struct {
	hand      hand.VisibleHand
	drawnTile *tile.Tile
}

func NewVisiblePlayer(handTiles []tile.Tile) (*VisiblePlayer, error) {
	if len(handTiles) != initHandSize {
		return nil, fmt.Errorf("invalid number of hand tiles: got %d, want %d", len(handTiles), initHandSize)
	}

	h, err := hand.NewVisibleHand(handTiles)
	if err != nil {
		return nil, err
	}

	return &VisiblePlayer{hand: *h, drawnTile: nil}, nil
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
