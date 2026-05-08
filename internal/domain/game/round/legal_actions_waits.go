package round

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/player/hand"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/round/service"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func canConcealedKanAfterRiichi(handBeforeKan *hand.VisibleHand, drawnTile tile.Tile, consumed [4]tile.Tile) bool {
	if !drawnTile.HasSameSymbol(consumed[0]) {
		return false
	}
	drawnTileID34 := drawnTile.RemoveRed().ID()
	tc34AfterKan := *handBeforeKan.ToTileCounts34()
	if tc34AfterKan[drawnTileID34] != 3 {
		return false
	}
	tc34AfterKan[drawnTileID34] = 0
	handAfterKan := hand.MustVisibleHand(tc34AfterKan.ToTiles())
	return service.WaitsFor(handBeforeKan) == service.WaitsFor(handAfterKan)
}
