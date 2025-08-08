package protocol

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type Adapter interface {
	DecodeMessages(msg []byte) ([]inbound.Event, error)
	EncodeResponse(ev outbound.Event) ([]byte, error)
}
