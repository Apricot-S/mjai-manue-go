package block

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/tile"
)

type Block interface {
	ToTiles() []tile.Tile
}
