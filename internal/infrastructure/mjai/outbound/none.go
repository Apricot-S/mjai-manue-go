package outbound

type None struct {
	Type string `json:"type"`
	Log  string `json:"log,omitempty"`
}

func NewNone(log string) *None {
	return &None{
		Type: "none",
		Log:  log,
	}
}

func (*None) outboundMessage() {}
