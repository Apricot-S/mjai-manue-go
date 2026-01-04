package hand

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

// For test only
func CodesToHand(codes []string) *VisibleHand {
	tiles := make([]tile.Tile, len(codes))
	for i, code := range codes {
		tiles[i] = *tile.MustTileFromCode(code)
	}

	h, err := NewVisibleHand(tiles)
	if err != nil {
		panic(err)
	}

	return h
}
