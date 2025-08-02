package event

type OutboundEvent interface {
	// isOutboundEvent is a marker method to distinguish outbound events.
	isOutboundEvent()
}
