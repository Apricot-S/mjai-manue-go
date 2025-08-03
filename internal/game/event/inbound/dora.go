package inbound

import "github.com/Apricot-S/mjai-manue-go/internal/base"

type Dora struct {
	DoraMarker base.Pai
}

func NewDora(doraMarker base.Pai) *Dora {
	event := &Dora{
		DoraMarker: doraMarker,
	}

	return event
}

func (n *Dora) isInboundEvent() {}
