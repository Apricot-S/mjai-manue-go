package outbound

type Message interface {
	outboundMessage()
}
