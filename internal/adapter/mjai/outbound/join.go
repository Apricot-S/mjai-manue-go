package outbound

type Join struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Room string `json:"room"`
}

func NewJoin(name string, room string) *Join {
	return &Join{
		Type: "join",
		Name: name,
		Room: room,
	}
}

func (*Join) outboundMessage() {}
