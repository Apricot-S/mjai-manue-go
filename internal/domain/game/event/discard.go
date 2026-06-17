package event

import (
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/seat"
	"github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"
)

type Discard struct {
	actor     seat.Seat
	tile      tile.Tile
	tsumogiri bool
}

func NewDiscard(actor seat.Seat, discardedTile tile.Tile, tsumogiri bool) *Discard {
	return &Discard{
		actor:     actor,
		tile:      discardedTile,
		tsumogiri: tsumogiri,
	}
}

func (*Discard) isEvent() {}

func (d *Discard) Actor() seat.Seat {
	return d.actor
}

func (d *Discard) Tile() tile.Tile {
	return d.tile
}

func (d *Discard) Tsumogiri() bool {
	return d.tsumogiri
}
