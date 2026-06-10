package player

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"

func isSwapCallTile(t tile.Tile, swapCallTiles []tile.Tile) bool {
	return tile.Tiles(swapCallTiles).ContainsSameSymbol(t)
}
