package inbound

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (*Error) inboundMessage() {}
