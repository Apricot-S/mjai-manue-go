package inbound

type Hello struct {
	Type            string `json:"type"`
	Protocol        string `json:"protocol"`
	ProtocolVersion int    `json:"protocol_version"`
}

func (*Hello) inboundMessage() {}
