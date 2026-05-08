package event

import "github.com/Apricot-S/mjai-manue-go/internal/domain/game/tile"

type Dora struct {
	indicator tile.Tile
}

func NewDora(indicator tile.Tile) *Dora {
	return &Dora{indicator: indicator}
}

func (*Dora) isEvent() {}

func (d *Dora) Indicator() tile.Tile {
	return d.indicator
}
