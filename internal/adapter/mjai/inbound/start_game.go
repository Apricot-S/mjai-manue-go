package inbound

type StartGame struct {
	Type  string   `json:"type"`
	ID    *int     `json:"id,omitempty"`
	Names []string `json:"names,omitempty"`
}

func (*StartGame) inboundMessage() {}
