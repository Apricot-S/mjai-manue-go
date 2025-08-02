package event

type InboundEvent interface {
	Type() InboundEventType
}
