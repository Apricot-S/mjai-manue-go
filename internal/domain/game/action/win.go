package action

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Win struct {
	actor       seat.Seat
	target      seat.Seat
	winningTile tile.Tile
}

func NewWin(actor, target seat.Seat, winningTile tile.Tile) (*Win, error) {
	if winningTile.IsUnknown() {
		return nil, fmt.Errorf("unknown tile is not allowed for Win action: %s", winningTile)
	}
	return &Win{actor: actor, target: target, winningTile: winningTile}, nil
}

func (*Win) isAction() {}

func (w *Win) Actor() seat.Seat {
	return w.actor
}

func (w *Win) Target() seat.Seat {
	return w.target
}

func (w *Win) WinningTile() tile.Tile {
	return w.winningTile
}
