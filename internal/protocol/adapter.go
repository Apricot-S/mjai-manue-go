package protocol

import (
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/inbound"
	"github.com/Apricot-S/mjai-manue-go/internal/game/event/outbound"
)

type MessageEventAdapter interface {
	MessageToEvent(msg []byte) (inbound.Event, error)
}

type EventMessageAdapter interface {
	EventToMessage(ev outbound.Event) ([]byte, error)
}
