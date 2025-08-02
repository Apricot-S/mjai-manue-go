package event

type OutboundEvent interface {
	Type() OutboundEventType
}
