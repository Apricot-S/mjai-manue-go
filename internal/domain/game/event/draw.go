package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Draw struct {
	actor     seat.Seat
	drawnTile tile.Tile
}

func NewDraw(actor seat.Seat, drawnTile tile.Tile) *Draw {
	return &Draw{actor, drawnTile}
}

func (*Draw) isEvent() {}

func (d *Draw) Actor() seat.Seat {
	return d.actor
}

func (d *Draw) DrawnTile() tile.Tile {
	return d.drawnTile
}
