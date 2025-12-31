package hand

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/model/tile"
)

const maxNumTilesInHand = 14
const maxCopies = 4

type Hand interface {
	ToTiles() []tile.Tile
	Draw(tile *tile.Tile) (Hand, error)
	// Discard(tile *tile.Tile) (Hand, error)
	// Call(meld meld.Meld) (Hand, error)
}
