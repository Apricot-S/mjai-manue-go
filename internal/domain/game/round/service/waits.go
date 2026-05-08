package service

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type WaitSet uint64

func (w WaitSet) Has(t tile.Tile) bool {
	return w&(WaitSet(1)<<t.RemoveRed().ID()) != 0
}

func WaitsFor(h *hand.VisibleHand) WaitSet {
	var waits WaitSet
	for id := range tile.NumTileType34 {
		waitTile := tile.MustTileFromID(id)
		handWithWait, err := h.Draw(waitTile)
		if err != nil {
			continue
		}
		if IsWinningForm(handWithWait) {
			waits |= WaitSet(1) << id
		}
	}
	return waits
}
