package outbound

type None struct {
	Type string `json:"type"`
}

func NewNone() *None {
	return &None{
		Type: "none",
	}
}

func (*None) outboundMessage() {}
