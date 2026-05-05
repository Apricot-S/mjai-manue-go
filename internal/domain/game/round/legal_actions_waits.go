package round

import (
	"maps"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

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
