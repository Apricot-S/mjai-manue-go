package action

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"

func hasUnknownTile(tiles []tile.Tile) bool {
	for _, t := range tiles {
		if t.IsUnknown() {
			return true
		}
	}
	return false
}
