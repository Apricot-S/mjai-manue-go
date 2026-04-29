package player

import (
	"slices"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

func isSwapCallTile(t tile.Tile, swapCallTiles []tile.Tile) bool {
	return slices.ContainsFunc(swapCallTiles, func(s tile.Tile) bool {
		return t.HasSameSymbol(&s)
	})
}
