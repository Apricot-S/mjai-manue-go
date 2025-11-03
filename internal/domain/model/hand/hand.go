package hand

import "github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"

type Hand struct {
	tiles []tile.Tile
}

func NewHand(tiles []tile.Tile) *Hand {
	return &Hand{tiles: tiles}
}

func (h *Hand) ToTiles() []tile.Tile {
	return h.tiles
}
