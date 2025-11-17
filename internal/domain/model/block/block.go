package block

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

type Block interface {
	ToTiles() []tile.Tile
}
