package inbound

type InboundEvent interface {
	// isInboundEvent is a marker method to distinguish inbound events.
	isInboundEvent()
}
