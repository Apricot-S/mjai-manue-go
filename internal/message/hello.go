package message

type Hello struct {
	Message
	Protocol        string `json:"protocol,omitempty"`
	ProtocolVersion int    `json:"protocol_version,omitempty"`
}

func NewHello(protocol string, protocolVersion int) *Hello {
	return &Hello{
		Message:         Message{Type: TypeHello},
		Protocol:        protocol,
		ProtocolVersion: protocolVersion,
	}
}
