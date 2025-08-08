package inbound

import (
	"fmt"

	"github.com/Apricot-S/mjai-manue-go/internal/base"
)

type Dora struct {
	DoraMarker base.Pai
}

func NewDora(doraMarker base.Pai) (*Dora, error) {
	event := &Dora{
		DoraMarker: doraMarker,
	}

	if event.DoraMarker.IsUnknown() {
		return nil, fmt.Errorf("doraMarker must not be unknown: %v", event)
	}

	return event, nil
}

func (d *Dora) isInboundEvent() {}
