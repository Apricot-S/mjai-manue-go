package inbound

type EndGame struct {
	Type   string `json:"type"`
	Scores []int  `json:"scores,omitempty"`
}

func (*EndGame) inboundMessage() {}
