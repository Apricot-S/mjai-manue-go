package outbound

type Event interface {
	// isOutboundEvent is a marker method to distinguish outbound events.
	isOutboundEvent()
}
