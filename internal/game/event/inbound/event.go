package inbound

type Event interface {
	// isInboundEvent is a marker method to distinguish inbound events.
	isInboundEvent()
}
